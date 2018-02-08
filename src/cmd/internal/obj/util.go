// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package obj

import (
	"bytes"
	"cmd/internal/objabi"
	"fmt"
	"strings"
)

const REG_NONE = 0

// Line returns a string containing the filename and line number for p
func (p *Prog) Line() string {
	return p.Ctxt.OutermostPos(p.Pos).Format(false, true)
}

// InnermostLineNumber returns a string containing the line number for the
// innermost inlined function (if any inlining) at p's position
func (p *Prog) InnermostLineNumber() string {
	pos := p.Ctxt.InnermostPos(p.Pos)
	if !pos.IsKnown() {
		return "?"
	}
	return fmt.Sprintf("%d", pos.Line())
}

// InnermostFilename returns a string containing the innermost
// (in inlining) filename at p's position
func (p *Prog) InnermostFilename() string {
	// TODO For now, this is only used for debugging output, and if we need more/better information, it might change.
	// An example of what we might want to see is the full stack of positions for inlined code, so we get some visibility into what is recorded there.
	pos := p.Ctxt.InnermostPos(p.Pos)
	if !pos.IsKnown() {
		return "<unknown file name>"
	}
	return pos.Filename()
}

var armCondCode = []string{
	".EQ",
	".NE",
	".CS",
	".CC",
	".MI",
	".PL",
	".VS",
	".VC",
	".HI",
	".LS",
	".GE",
	".LT",
	".GT",
	".LE",
	"",
	".NV",
}

/* ARM scond byte */
const (
	C_SCOND     = (1 << 4) - 1
	C_SBIT      = 1 << 4
	C_PBIT      = 1 << 5
	C_WBIT      = 1 << 6
	C_FBIT      = 1 << 7
	C_UBIT      = 1 << 7
	C_SCOND_XOR = 14
)

// CConv formats ARM condition codes.
func CConv(s uint8) string {
	if s == 0 {
		return ""
	}
	sc := armCondCode[(s&C_SCOND)^C_SCOND_XOR]
	if s&C_SBIT != 0 {
		sc += ".S"
	}
	if s&C_PBIT != 0 {
		sc += ".P"
	}
	if s&C_WBIT != 0 {
		sc += ".W"
	}
	if s&C_UBIT != 0 { /* ambiguous with FBIT */
		sc += ".U"
	}
	return sc
}

func (p *Prog) String() string {
	if p == nil {
		return "<nil Prog>"
	}
	if p.Ctxt == nil {
		return "<Prog without ctxt>"
	}
	return fmt.Sprintf("%.5d (%v)\t%s", p.Pc, p.Line(), p.InstructionString())
}

// InstructionString returns a string representation of the instruction without preceding
// program counter or file and line number.
func (p *Prog) InstructionString() string {
	if p == nil {
		return "<nil Prog>"
	}

	if p.Ctxt == nil {
		return "<Prog without ctxt>"
	}

	sc := CConv(p.Scond)

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%v%s", p.As, sc)
	sep := "\t"

	if p.From.Type != TYPE_NONE {
		fmt.Fprintf(&buf, "%s%v", sep, Dconv(p, &p.From))
		sep = ", "
	}
	if p.Reg != REG_NONE {
		// Should not happen but might as well show it if it does.
		fmt.Fprintf(&buf, "%s%v", sep, Rconv(int(p.Reg)))
		sep = ", "
	}
	for i := range p.RestArgs {
		fmt.Fprintf(&buf, "%s%v", sep, Dconv(p, &p.RestArgs[i]))
		sep = ", "
	}

	if p.As == ATEXT {
		// If there are attributes, print them. Otherwise, skip the comma.
		// In short, print one of these two:
		// TEXT	foo(SB), DUPOK|NOSPLIT, $0
		// TEXT	foo(SB), $0
		s := p.From.Sym.Attribute.TextAttrString()
		if s != "" {
			fmt.Fprintf(&buf, "%s%s", sep, s)
			sep = ", "
		}
	}
	if p.To.Type != TYPE_NONE {
		fmt.Fprintf(&buf, "%s%v", sep, Dconv(p, &p.To))
	}
	if p.RegTo2 != REG_NONE {
		fmt.Fprintf(&buf, "%s%v", sep, Rconv(int(p.RegTo2)))
	}
	return buf.String()
}

func (ctxt *Link) NewProg() *Prog {
	p := new(Prog)
	p.Ctxt = ctxt
	return p
}

func (ctxt *Link) CanReuseProgs() bool {
	return !ctxt.Debugasm
}

func (ctxt *Link) Dconv(a *Addr) string {
	return Dconv(nil, a)
}

func Dconv(p *Prog, a *Addr) string {
	var str string

	switch a.Type {
	default:
		str = fmt.Sprintf("type=%d", a.Type)

	case TYPE_NONE:
		str = ""
		if a.Name != NAME_NONE || a.Reg != 0 || a.Sym != nil {
			str = fmt.Sprintf("%v(%v)(NONE)", Mconv(a), Rconv(int(a.Reg)))
		}

	case TYPE_REG:
		// TODO(rsc): This special case is for x86 instructions like
		//	PINSRQ	CX,$1,X6
		// where the $1 is included in the p->to Addr.
		// Move into a new field.
		if a.Offset != 0 && (a.Reg < RBaseARM64 || a.Reg >= RBaseMIPS) {
			str = fmt.Sprintf("$%d,%v", a.Offset, Rconv(int(a.Reg)))
			break
		}

		str = Rconv(int(a.Reg))
		if a.Name != NAME_NONE || a.Sym != nil {
			str = fmt.Sprintf("%v(%v)(REG)", Mconv(a), Rconv(int(a.Reg)))
		}
		if (RBaseARM64+1<<10+1<<9) /* arm64.REG_ELEM */ <= a.Reg &&
			a.Reg < (RBaseARM64+1<<11) /* arm64.REG_ELEM_END */ {
			str += fmt.Sprintf("[%d]", a.Index)
		}

	case TYPE_BRANCH:
		if a.Sym != nil {
			str = fmt.Sprintf("%s(SB)", a.Sym.Name)
		} else if p != nil && p.Pcond != nil {
			str = fmt.Sprint(p.Pcond.Pc)
		} else if a.Val != nil {
			str = fmt.Sprint(a.Val.(*Prog).Pc)
		} else {
			str = fmt.Sprintf("%d(PC)", a.Offset)
		}

	case TYPE_INDIR:
		str = fmt.Sprintf("*%s", Mconv(a))

	case TYPE_MEM:
		str = Mconv(a)
		if a.Index != REG_NONE {
			str += fmt.Sprintf("(%v*%d)", Rconv(int(a.Index)), int(a.Scale))
		}

	case TYPE_CONST:
		if a.Reg != 0 {
			str = fmt.Sprintf("$%v(%v)", Mconv(a), Rconv(int(a.Reg)))
		} else {
			str = fmt.Sprintf("$%v", Mconv(a))
		}

	case TYPE_TEXTSIZE:
		if a.Val.(int32) == objabi.ArgsSizeUnknown {
			str = fmt.Sprintf("$%d", a.Offset)
		} else {
			str = fmt.Sprintf("$%d-%d", a.Offset, a.Val.(int32))
		}

	case TYPE_FCONST:
		str = fmt.Sprintf("%.17g", a.Val.(float64))
		// Make sure 1 prints as 1.0
		if !strings.ContainsAny(str, ".e") {
			str += ".0"
		}
		str = fmt.Sprintf("$(%s)", str)

	case TYPE_SCONST:
		str = fmt.Sprintf("$%q", a.Val.(string))

	case TYPE_ADDR:
		str = fmt.Sprintf("$%s", Mconv(a))

	case TYPE_SHIFT:
		v := int(a.Offset)
		ops := "<<>>->@>"
		switch objabi.GOARCH {
		case "arm":
			op := ops[((v>>5)&3)<<1:]
			if v&(1<<4) != 0 {
				str = fmt.Sprintf("R%d%c%cR%d", v&15, op[0], op[1], (v>>8)&15)
			} else {
				str = fmt.Sprintf("R%d%c%c%d", v&15, op[0], op[1], (v>>7)&31)
			}
			if a.Reg != 0 {
				str += fmt.Sprintf("(%v)", Rconv(int(a.Reg)))
			}
		case "arm64":
			op := ops[((v>>22)&3)<<1:]
			str = fmt.Sprintf("R%d%c%c%d", (v>>16)&31, op[0], op[1], (v>>10)&63)
		default:
			panic("TYPE_SHIFT is not supported on " + objabi.GOARCH)
		}

	case TYPE_REGREG:
		str = fmt.Sprintf("(%v, %v)", Rconv(int(a.Reg)), Rconv(int(a.Offset)))

	case TYPE_REGREG2:
		str = fmt.Sprintf("%v, %v", Rconv(int(a.Offset)), Rconv(int(a.Reg)))

	case TYPE_REGLIST:
		str = RLconv(a.Offset)
	}

	return str
}

func Mconv(a *Addr) string {
	var str string

	switch a.Name {
	default:
		str = fmt.Sprintf("name=%d", a.Name)

	case NAME_NONE:
		switch {
		case a.Reg == REG_NONE:
			str = fmt.Sprint(a.Offset)
		case a.Offset == 0:
			str = fmt.Sprintf("(%v)", Rconv(int(a.Reg)))
		case a.Offset != 0:
			str = fmt.Sprintf("%d(%v)", a.Offset, Rconv(int(a.Reg)))
		}

		// Note: a.Reg == REG_NONE encodes the default base register for the NAME_ type.
	case NAME_EXTERN:
		reg := "SB"
		if a.Reg != REG_NONE {
			reg = Rconv(int(a.Reg))
		}
		if a.Sym != nil {
			str = fmt.Sprintf("%s%s(%s)", a.Sym.Name, offConv(a.Offset), reg)
		} else {
			str = fmt.Sprintf("%s(%s)", offConv(a.Offset), reg)
		}

	case NAME_GOTREF:
		reg := "SB"
		if a.Reg != REG_NONE {
			reg = Rconv(int(a.Reg))
		}
		if a.Sym != nil {
			str = fmt.Sprintf("%s%s@GOT(%s)", a.Sym.Name, offConv(a.Offset), reg)
		} else {
			str = fmt.Sprintf("%s@GOT(%s)", offConv(a.Offset), reg)
		}

	case NAME_STATIC:
		reg := "SB"
		if a.Reg != REG_NONE {
			reg = Rconv(int(a.Reg))
		}
		if a.Sym != nil {
			str = fmt.Sprintf("%s<>%s(%s)", a.Sym.Name, offConv(a.Offset), reg)
		} else {
			str = fmt.Sprintf("<>%s(%s)", offConv(a.Offset), reg)
		}

	case NAME_AUTO:
		reg := "SP"
		if a.Reg != REG_NONE {
			reg = Rconv(int(a.Reg))
		}
		if a.Sym != nil {
			str = fmt.Sprintf("%s%s(%s)", a.Sym.Name, offConv(a.Offset), reg)
		} else {
			str = fmt.Sprintf("%s(%s)", offConv(a.Offset), reg)
		}

	case NAME_PARAM:
		reg := "FP"
		if a.Reg != REG_NONE {
			reg = Rconv(int(a.Reg))
		}
		if a.Sym != nil {
			str = fmt.Sprintf("%s%s(%s)", a.Sym.Name, offConv(a.Offset), reg)
		} else {
			str = fmt.Sprintf("%s(%s)", offConv(a.Offset), reg)
		}
	}
	return str
}

func offConv(off int64) string {
	if off == 0 {
		return ""
	}
	return fmt.Sprintf("%+d", off)
}

type regSet struct {
	lo    int
	hi    int
	Rconv func(int) string
}

// Few enough architectures that a linear scan is fastest.
// Not even worth sorting.
var regSpace []regSet

/*
	Each architecture defines a register space as a unique
	integer range.
	Here is the list of architectures and the base of their register spaces.
*/

const (
	// Because of masking operations in the encodings, each register
	// space should start at 0 modulo some power of 2.
	RBase386   = 1 * 1024
	RBaseAMD64 = 2 * 1024
	RBaseARM   = 3 * 1024
	RBasePPC64 = 4 * 1024  // range [4k, 8k)
	RBaseARM64 = 8 * 1024  // range [8k, 13k)
	RBaseMIPS  = 13 * 1024 // range [13k, 14k)
	RBaseS390X = 14 * 1024 // range [14k, 15k)
	RBaseWASM  = 16 * 1024
)

// RegisterRegister binds a pretty-printer (Rconv) for register
// numbers to a given register number range. Lo is inclusive,
// hi exclusive (valid registers are lo through hi-1).
func RegisterRegister(lo, hi int, Rconv func(int) string) {
	regSpace = append(regSpace, regSet{lo, hi, Rconv})
}

func Rconv(reg int) string {
	if reg == REG_NONE {
		return "NONE"
	}
	for i := range regSpace {
		rs := &regSpace[i]
		if rs.lo <= reg && reg < rs.hi {
			return rs.Rconv(reg)
		}
	}
	return fmt.Sprintf("R???%d", reg)
}

type regListSet struct {
	lo     int64
	hi     int64
	RLconv func(int64) string
}

var regListSpace []regListSet

// Each architecture is allotted a distinct subspace: [Lo, Hi) for declaring its
// arch-specific register list numbers.
const (
	RegListARMLo = 0
	RegListARMHi = 1 << 16

	// arm64 uses the 60th bit to differentiate from other archs
	RegListARM64Lo = 1 << 60
	RegListARM64Hi = 1<<61 - 1
)

// RegisterRegisterList binds a pretty-printer (RLconv) for register list
// numbers to a given register list number range. Lo is inclusive,
// hi exclusive (valid register list are lo through hi-1).
func RegisterRegisterList(lo, hi int64, rlconv func(int64) string) {
	regListSpace = append(regListSpace, regListSet{lo, hi, rlconv})
}

func RLconv(list int64) string {
	for i := range regListSpace {
		rls := &regListSpace[i]
		if rls.lo <= list && list < rls.hi {
			return rls.RLconv(list)
		}
	}
	return fmt.Sprintf("RL???%d", list)
}

type opSet struct {
	lo    As
	names []string
}

// Not even worth sorting
var aSpace []opSet

// RegisterOpcode binds a list of instruction names
// to a given instruction number range.
func RegisterOpcode(lo As, Anames []string) {
	if len(Anames) > AllowedOpCodes {
		panic(fmt.Sprintf("too many instructions, have %d max %d", len(Anames), AllowedOpCodes))
	}
	aSpace = append(aSpace, opSet{lo, Anames})
}

func (a As) String() string {
	if 0 <= a && int(a) < len(Anames) {
		return Anames[a]
	}
	for i := range aSpace {
		as := &aSpace[i]
		if as.lo <= a && int(a-as.lo) < len(as.names) {
			return as.names[a-as.lo]
		}
	}
	return fmt.Sprintf("A???%d", a)
}

var Anames = []string{
	"XXX",
	"CALL",
	"DUFFCOPY",
	"DUFFZERO",
	"END",
	"FUNCDATA",
	"JMP",
	"NOP",
	"PCDATA",
	"RET",
	"TEXT",
	"UNDEF",
}

func Bool2int(b bool) int {
	// The compiler currently only optimizes this form.
	// See issue 6011.
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
