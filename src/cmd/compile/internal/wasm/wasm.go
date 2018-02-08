// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"cmd/compile/internal/gc"
	"cmd/compile/internal/ssa"
	"cmd/compile/internal/types"
	"cmd/internal/obj"
	"cmd/internal/obj/wasm"
)

func Init(arch *gc.Arch) {
	arch.LinkArch = &wasm.Linkwasm
	arch.REGSP = wasm.REG_SP
	arch.MAXWIDTH = 1 << 50

	arch.ZeroRange = zeroRange
	arch.ZeroAuto = zeroAuto
	arch.Ginsnop = ginsnop

	arch.SSAMarkMoves = ssaMarkMoves
	arch.SSAGenValue = ssaGenValue
	arch.SSAGenBlock = ssaGenBlock
}

func zeroRange(pp *gc.Progs, p *obj.Prog, off, cnt int64, state *uint32) *obj.Prog {
	if cnt == 0 {
		return p
	}
	if cnt%8 != 0 {
		gc.Fatalf("zerorange count not a multiple of widthptr %d", cnt)
	}

	for i := int64(0); i < cnt; i += 8 {
		p = pp.Appendpp(p, wasm.AGet, obj.TYPE_REG, wasm.REG_SP, 0, 0, 0, 0)
		p = pp.Appendpp(p, wasm.AI64Const, obj.TYPE_CONST, 0, 0, 0, 0, 0)
		p = pp.Appendpp(p, wasm.AI64Store, 0, 0, 0, obj.TYPE_CONST, 0, off+i)
	}

	return p
}

func zeroAuto(pp *gc.Progs, n *gc.Node) {
	sym := n.Sym.Linksym()
	size := n.Type.Size()
	for i := int64(0); i < size; i += 8 {
		p := pp.Prog(wasm.AGet)
		p.From = obj.Addr{Type: obj.TYPE_REG, Reg: wasm.REG_SP}

		p = pp.Prog(wasm.AI64Const)
		p.From = obj.Addr{Type: obj.TYPE_CONST, Offset: 0}

		p = pp.Prog(wasm.AI64Store)
		p.To = obj.Addr{Type: obj.TYPE_MEM, Name: obj.NAME_AUTO, Offset: n.Xoffset + i, Sym: sym}
	}
}

func ginsnop(pp *gc.Progs) {
	pp.Prog(wasm.ANop)
}

func ssaMarkMoves(s *gc.SSAGenState, b *ssa.Block) {
}

func ssaGenBlock(s *gc.SSAGenState, b, next *ssa.Block) {
	goToBlock := func(block *ssa.Block, canFallthrough bool) {
		if canFallthrough && block == next {
			return
		}
		p := s.Prog(obj.AJMP)
		p.To.Type = obj.TYPE_BRANCH
		s.Branches = append(s.Branches, gc.Branch{P: p, B: block})
	}

	switch b.Kind {
	case ssa.BlockPlain:
		goToBlock(b.Succs[0].Block(), true)

	case ssa.BlockIf:
		getReg32(s, b.Control)
		s.Prog(wasm.AI32Eqz)
		s.Prog(wasm.AIf)
		goToBlock(b.Succs[1].Block(), false)
		s.Prog(wasm.AEnd)
		goToBlock(b.Succs[0].Block(), true)

	case ssa.BlockRet:
		p2 := s.Prog(wasm.ABr)
		p2.To = obj.Addr{Type: obj.TYPE_CONST, Offset: -2}

	case ssa.BlockRetJmp:
		p := s.Prog(obj.AJMP)
		p.To.Type = obj.TYPE_MEM
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = b.Aux.(*obj.LSym)

	case ssa.BlockExit:
		s.Prog(wasm.AUnreachable)

	case ssa.BlockDefer:
		s.Prog(wasm.ABlock)
		p := s.Prog(wasm.AGet)
		p.From = obj.Addr{Type: obj.TYPE_REG, Reg: wasm.REG_RET0}
		s.Prog(wasm.AI64Eqz)
		s.Prog(wasm.ABrIf)
		goToBlock(b.Succs[1].Block(), false)
		s.Prog(wasm.AEnd)
		goToBlock(b.Succs[0].Block(), true)

	default:
		panic("unexpected block")
	}
	s.Prog(wasm.AEnd)
}

func ssaGenValue(s *gc.SSAGenState, v *ssa.Value) {
	switch v.Op {
	case ssa.OpWasmLoweredStaticCall, ssa.OpWasmLoweredClosureCall, ssa.OpWasmLoweredInterCall:
		if v.Aux == gc.Deferreturn {
			// start new block before call to deferreturn so it can be called again via jmpdefer
			s.Prog(wasm.AEnd)
		}
		if v.Op == ssa.OpWasmLoweredClosureCall {
			getReg64(s, v.Args[1])
			setReg(s, wasm.REG_CTX)
		}
		p := s.Call(v)
		p.Pos = v.Pos
		s.Prog(wasm.AEnd)

	case ssa.OpWasmLoweredConvert:
		if v.Args[0].Reg() != v.Reg() {
			v.Fatalf("OpWasmLoweredConvert should be a no-op")
		}

	case ssa.OpWasmLoweredMove:
		getReg32(s, v.Args[0])
		getReg32(s, v.Args[1])
		i32Const(s, int32(v.AuxInt))
		p := s.Prog(wasm.ACall)
		p.To = obj.Addr{Type: obj.TYPE_MEM, Sym: gc.WasmMove}

	case ssa.OpWasmLoweredZero:
		getReg32(s, v.Args[0])
		i32Const(s, int32(v.AuxInt))
		p := s.Prog(wasm.ACall)
		p.To = obj.Addr{Type: obj.TYPE_MEM, Sym: gc.WasmZero}

	case ssa.OpWasmLoweredNilCheck:
		getReg64(s, v.Args[0])
		s.Prog(wasm.AI64Eqz)
		s.Prog(wasm.AI32Eqz)
		s.Prog(wasm.ABrIf)
		p := s.Prog(obj.ACALL)
		p.To = obj.Addr{Type: obj.TYPE_MEM, Sym: gc.SigPanic}
		s.Prog(wasm.AEnd)
		if gc.Debug_checknil != 0 && v.Pos.Line() > 1 { // v.Pos.Line()==1 in generated wrappers
			gc.Warnl(v.Pos, "generated nil check")
		}

	case ssa.OpWasmI64Store8, ssa.OpWasmI64Store16, ssa.OpWasmI64Store32, ssa.OpWasmI64Store, ssa.OpWasmF32Store, ssa.OpWasmF64Store:
		getReg32(s, v.Args[0])
		getReg64(s, v.Args[1])
		if v.Op == ssa.OpWasmF32Store {
			s.Prog(wasm.AF32DemoteF64)
		}
		p := s.Prog(v.Op.Asm())
		p.To = obj.Addr{Type: obj.TYPE_MEM, Offset: v.AuxInt}

	case ssa.OpWasmI64Store8Const, ssa.OpWasmI64Store16Const, ssa.OpWasmI64Store32Const, ssa.OpWasmI64StoreConst:
		valOff := v.AuxValAndOff()
		getReg32(s, v.Args[0])
		i64Const(s, valOff.Val())
		p := s.Prog(v.Op.Asm())
		p.To = obj.Addr{Type: obj.TYPE_MEM, Offset: valOff.Off()}

	case ssa.OpStoreReg:
		getReg(s, wasm.REG_SP)
		getReg64(s, v.Args[0])
		if v.Type.Etype == types.TFLOAT32 {
			s.Prog(wasm.AF32DemoteF64)
		}
		p := s.Prog(storeOp(v.Type))
		gc.AddrAuto(&p.To, v)

	default:
		if !v.HasReg() {
			return
		}
		ssaGenValueOnStack(s, v)
		setReg(s, v.Reg())
	}
}

func ssaGenValueOnStack(s *gc.SSAGenState, v *ssa.Value) {
	switch v.Op {
	case ssa.OpWasmLoweredGetClosurePtr:
		getReg(s, wasm.REG_CTX)

	case ssa.OpWasmLoweredGetCallerPC:
		p := s.Prog(wasm.AI64Load)
		p.From = obj.Addr{Type: obj.TYPE_MEM, Name: obj.NAME_PARAM, Offset: -8} // PC is stored 8 bytes below first parameter.

	case ssa.OpWasmLoweredGetCallerSP:
		p := s.Prog(wasm.AGet)
		p.From = obj.Addr{Type: obj.TYPE_REG, Reg: wasm.REG_SP}
		s.Prog(wasm.AI64ExtendUI32)
		p = s.Prog(wasm.AI64Const)
		p.From = obj.Addr{Type: obj.TYPE_ADDR, Name: obj.NAME_PARAM, Offset: 0}
		s.Prog(wasm.AI64Add)

	case ssa.OpWasmLoweredAddr:
		p := s.Prog(wasm.AAddrOf)
		p.From = obj.Addr{Type: obj.TYPE_MEM, Reg: v.Args[0].Reg()}
		gc.AddAux(&p.From, v)

	case ssa.OpWasmLoweredRound32F:
		getReg64(s, v.Args[0])
		s.Prog(wasm.AF32DemoteF64)
		s.Prog(wasm.AF64PromoteF32)

	case ssa.OpWasmSelect:
		getReg64(s, v.Args[0])
		getReg64(s, v.Args[1])
		getReg64(s, v.Args[2])
		s.Prog(wasm.AI32WrapI64)
		s.Prog(v.Op.Asm())

	case ssa.OpWasmI64AddConst:
		getReg64(s, v.Args[0])
		i64Const(s, v.AuxInt)
		s.Prog(v.Op.Asm())

	case ssa.OpWasmI64Const:
		i64Const(s, v.AuxInt)

	case ssa.OpWasmF64Const:
		f64Const(s, v.AuxFloat())

	case ssa.OpWasmI64Load8U, ssa.OpWasmI64Load8S, ssa.OpWasmI64Load16U, ssa.OpWasmI64Load16S, ssa.OpWasmI64Load32U, ssa.OpWasmI64Load32S, ssa.OpWasmI64Load, ssa.OpWasmF32Load, ssa.OpWasmF64Load:
		getReg32(s, v.Args[0])
		p := s.Prog(v.Op.Asm())
		p.From = obj.Addr{Type: obj.TYPE_CONST, Offset: v.AuxInt}
		if v.Op == ssa.OpWasmF32Load {
			s.Prog(wasm.AF64PromoteF32)
		}

	case ssa.OpWasmI64Eqz:
		getReg64(s, v.Args[0])
		s.Prog(v.Op.Asm())
		s.Prog(wasm.AI64ExtendUI32)

	case ssa.OpWasmI64Eq, ssa.OpWasmI64Ne, ssa.OpWasmI64LtS, ssa.OpWasmI64LtU, ssa.OpWasmI64GtS, ssa.OpWasmI64GtU, ssa.OpWasmI64LeS, ssa.OpWasmI64LeU, ssa.OpWasmI64GeS, ssa.OpWasmI64GeU, ssa.OpWasmF64Eq, ssa.OpWasmF64Ne, ssa.OpWasmF64Lt, ssa.OpWasmF64Gt, ssa.OpWasmF64Le, ssa.OpWasmF64Ge:
		getReg64(s, v.Args[0])
		getReg64(s, v.Args[1])
		s.Prog(v.Op.Asm())
		s.Prog(wasm.AI64ExtendUI32)

	case ssa.OpWasmI64Add, ssa.OpWasmI64Sub, ssa.OpWasmI64Mul, ssa.OpWasmI64DivU, ssa.OpWasmI64RemS, ssa.OpWasmI64RemU, ssa.OpWasmI64And, ssa.OpWasmI64Or, ssa.OpWasmI64Xor, ssa.OpWasmI64Shl, ssa.OpWasmI64ShrS, ssa.OpWasmI64ShrU, ssa.OpWasmF64Add, ssa.OpWasmF64Sub, ssa.OpWasmF64Mul, ssa.OpWasmF64Div:
		getReg64(s, v.Args[0])
		getReg64(s, v.Args[1])
		s.Prog(v.Op.Asm())

	case ssa.OpWasmI64DivS:
		getReg64(s, v.Args[0])
		getReg64(s, v.Args[1])
		p := s.Prog(wasm.ACall)
		p.To = obj.Addr{Type: obj.TYPE_MEM, Sym: gc.WasmDiv}

	case ssa.OpWasmI64TruncSF64:
		getReg64(s, v.Args[0])
		p := s.Prog(wasm.ACall)
		p.To = obj.Addr{Type: obj.TYPE_MEM, Sym: gc.WasmTrunc}

	case ssa.OpWasmF64Neg, ssa.OpWasmF64ConvertSI64:
		getReg64(s, v.Args[0])
		s.Prog(v.Op.Asm())

	case ssa.OpLoadReg:
		p := s.Prog(loadOp(v.Type))
		gc.AddrAuto(&p.From, v.Args[0])
		if v.Type.Etype == types.TFLOAT32 {
			s.Prog(wasm.AF64PromoteF32)
		}

	case ssa.OpCopy:
		if v.Args[0].Type.IsMemory() {
			return
		}
		getReg64(s, v.Args[0])

	default:
		v.Fatalf("unexpected op: %s", v.Op)

	}
}

func getReg32(s *gc.SSAGenState, v *ssa.Value) {
	if !v.HasReg() {
		ssaGenValueOnStack(s, v)
		s.Prog(wasm.AI32WrapI64)
		return
	}

	reg := v.Reg()
	getReg(s, reg)
	if reg != wasm.REG_SP {
		s.Prog(wasm.AI32WrapI64)
	}
}

func getReg64(s *gc.SSAGenState, v *ssa.Value) {
	if !v.HasReg() {
		ssaGenValueOnStack(s, v)
		return
	}

	reg := v.Reg()
	getReg(s, reg)
	if reg == wasm.REG_SP {
		s.Prog(wasm.AI64ExtendUI32)
	}
}

func i32Const(s *gc.SSAGenState, val int32) {
	p := s.Prog(wasm.AI32Const)
	p.From = obj.Addr{Type: obj.TYPE_CONST, Offset: int64(val)}
}

func i64Const(s *gc.SSAGenState, val int64) {
	p := s.Prog(wasm.AI64Const)
	p.From = obj.Addr{Type: obj.TYPE_CONST, Offset: val}
}

func f64Const(s *gc.SSAGenState, val float64) {
	p := s.Prog(wasm.AF64Const)
	p.From = obj.Addr{Type: obj.TYPE_FCONST, Val: val}
}

func getReg(s *gc.SSAGenState, reg int16) {
	p := s.Prog(wasm.AGet)
	p.From = obj.Addr{Type: obj.TYPE_REG, Reg: reg}
}

func setReg(s *gc.SSAGenState, reg int16) {
	p := s.Prog(wasm.ASet)
	p.To = obj.Addr{Type: obj.TYPE_REG, Reg: reg}
}

func loadOp(t *types.Type) obj.As {
	if t.IsFloat() {
		switch t.Size() {
		case 4:
			return wasm.AF32Load
		case 8:
			return wasm.AF64Load
		default:
			panic("unexpected")
		}
	}

	switch t.Size() {
	case 1:
		if t.IsSigned() {
			return wasm.AI64Load8S
		}
		return wasm.AI64Load8U
	case 2:
		if t.IsSigned() {
			return wasm.AI64Load16S
		}
		return wasm.AI64Load16U
	case 4:
		if t.IsSigned() {
			return wasm.AI64Load32S
		}
		return wasm.AI64Load32U
	case 8:
		return wasm.AI64Load
	default:
		panic("unexpected")
	}
}

func storeOp(t *types.Type) obj.As {
	if t.IsFloat() {
		switch t.Size() {
		case 4:
			return wasm.AF32Store
		case 8:
			return wasm.AF64Store
		default:
			panic("unexpected")
		}
	}

	switch t.Size() {
	case 1:
		return wasm.AI64Store8
	case 2:
		return wasm.AI64Store16
	case 4:
		return wasm.AI64Store32
	case 8:
		return wasm.AI64Store
	default:
		panic("unexpected")
	}
}
