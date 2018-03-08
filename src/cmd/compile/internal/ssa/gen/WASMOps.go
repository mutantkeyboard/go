// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "strings"

var regNamesWasm = []string{
	"I0",
	"I1",
	"I2",
	"I3",
	"I4",
	"I5",
	"I6",
	"I7",
	"I8",
	"I9",
	"I10",
	"I11",
	"I12",
	"I13",
	"I14",
	"I15",

	"F0",
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"F8",
	"F9",
	"F10",
	"F11",
	"F12",
	"F13",
	"F14",
	"F15",

	"SP",
	"g",

	// pseudo-registers
	"SB",
}

func init() {
	// Make map from reg names to reg integers.
	if len(regNamesWasm) > 64 {
		panic("too many registers")
	}
	num := map[string]int{}
	for i, name := range regNamesWasm {
		num[name] = i
	}
	buildReg := func(s string) regMask {
		m := regMask(0)
		for _, r := range strings.Split(s, " ") {
			if n, ok := num[r]; ok {
				m |= regMask(1) << uint(n)
				continue
			}
			panic("register " + r + " not found")
		}
		return m
	}

	var (
		gp         = buildReg("I0 I1 I2 I3 I4 I5 I6 I7 I8 I9 I10 I11 I12 I13 I14 I15")
		fp         = buildReg("F0 F1 F2 F3 F4 F5 F6 F7 F8 F9 F10 F11 F12 F13 F14 F15")
		gpsp       = gp | buildReg("SP")
		gpspsb     = gpsp | buildReg("SB")
		callerSave = gp | fp
	)

	// Common regInfo
	var (
		gp01         = regInfo{inputs: nil, outputs: []regMask{gp}}
		gp11         = regInfo{inputs: []regMask{gpsp}, outputs: []regMask{gp}}
		gp21         = regInfo{inputs: []regMask{gpsp, gpsp}, outputs: []regMask{gp}}
		gp31         = regInfo{inputs: []regMask{gpsp, gpsp, gpsp}, outputs: []regMask{gp}}
		fp01         = regInfo{inputs: nil, outputs: []regMask{fp}}
		fp11         = regInfo{inputs: []regMask{fp}, outputs: []regMask{fp}}
		fp21         = regInfo{inputs: []regMask{fp, fp}, outputs: []regMask{fp}}
		fp21gp       = regInfo{inputs: []regMask{fp, fp}, outputs: []regMask{gp}}
		gpload       = regInfo{inputs: []regMask{gpspsb, 0}, outputs: []regMask{gp}}
		gpstore      = regInfo{inputs: []regMask{gpspsb, gpsp, 0}}
		gpstoreconst = regInfo{inputs: []regMask{gpspsb, 0}}
		fpload       = regInfo{inputs: []regMask{gpspsb, 0}, outputs: []regMask{fp}}
		fpstore      = regInfo{inputs: []regMask{gpspsb, fp, 0}}
		// fpstoreconst = regInfo{inputs: []regMask{fp, 0}}
	)

	var WasmOps = []opData{
		{name: "LoweredStaticCall", argLength: 1, reg: regInfo{clobbers: callerSave}, aux: "SymOff", call: true, symEffect: "None"},                         // call static function aux.(*obj.LSym).  arg0=mem, auxint=argsize, returns mem
		{name: "LoweredClosureCall", argLength: 3, reg: regInfo{inputs: []regMask{gp, buildReg("I15"), 0}, clobbers: callerSave}, aux: "Int64", call: true}, // call function via closure.  arg0=codeptr, arg1=closure, arg2=mem, auxint=argsize, returns mem
		{name: "LoweredInterCall", argLength: 2, reg: regInfo{inputs: []regMask{gp}, clobbers: callerSave}, aux: "Int64", call: true},                       // call fn by pointer.  arg0=codeptr, arg1=mem, auxint=argsize, returns mem

		{name: "LoweredAddr", argLength: 1, reg: gp11, aux: "SymOff", rematerializeable: true, symEffect: "Addr"},
		{name: "LoweredMove", argLength: 3, reg: regInfo{inputs: []regMask{gp, gp}}, aux: "Int64"},
		{name: "LoweredZero", argLength: 2, reg: regInfo{inputs: []regMask{gp}}, aux: "Int64"},
		{name: "LoweredConvert", argLength: 2, reg: regInfo{inputs: []regMask{gp}, outputs: []regMask{gp}}},
		{name: "LoweredGetClosurePtr", reg: gp01},
		{name: "LoweredGetCallerPC", reg: gp01, rematerializeable: true},
		{name: "LoweredGetCallerSP", reg: gp01, rematerializeable: true},
		{name: "LoweredNilCheck", argLength: 2, reg: regInfo{inputs: []regMask{gp}}, nilCheck: true, faultOnNilArg0: true},
		{name: "LoweredWB", argLength: 3, reg: regInfo{inputs: []regMask{gp, gp}}, aux: "Sym", symEffect: "None"},

		{name: "Select", asm: "Select", argLength: 3, reg: gp31},

		{name: "I64Load8U", asm: "I64Load8U", argLength: 2, reg: gpload, aux: "SymOff", typ: "UInt8", symEffect: "Read"},
		{name: "I64Load8S", asm: "I64Load8S", argLength: 2, reg: gpload, aux: "SymOff", typ: "Int8", symEffect: "Read"},
		{name: "I64Load16U", asm: "I64Load16U", argLength: 2, reg: gpload, aux: "SymOff", typ: "UInt16", symEffect: "Read"},
		{name: "I64Load16S", asm: "I64Load16S", argLength: 2, reg: gpload, aux: "SymOff", typ: "Int16", symEffect: "Read"},
		{name: "I64Load32U", asm: "I64Load32U", argLength: 2, reg: gpload, aux: "SymOff", typ: "UInt32", symEffect: "Read"},
		{name: "I64Load32S", asm: "I64Load32S", argLength: 2, reg: gpload, aux: "SymOff", typ: "Int32", symEffect: "Read"},
		{name: "I64Load", asm: "I64Load", argLength: 2, reg: gpload, aux: "SymOff", typ: "UInt64", symEffect: "Read"},
		{name: "I64Store8", asm: "I64Store8", argLength: 3, reg: gpstore, aux: "SymOff", typ: "Mem", symEffect: "Write"},
		{name: "I64Store8Const", asm: "I64Store8", argLength: 2, reg: gpstoreconst, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"},
		{name: "I64Store16", asm: "I64Store16", argLength: 3, reg: gpstore, aux: "SymOff", typ: "Mem", symEffect: "Write"},
		{name: "I64Store16Const", asm: "I64Store16", argLength: 2, reg: gpstoreconst, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"},
		{name: "I64Store32", asm: "I64Store32", argLength: 3, reg: gpstore, aux: "SymOff", typ: "Mem", symEffect: "Write"},
		{name: "I64Store32Const", asm: "I64Store32", argLength: 2, reg: gpstoreconst, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"},
		{name: "I64Store", asm: "I64Store", argLength: 3, reg: gpstore, aux: "SymOff", typ: "Mem", symEffect: "Write"},
		{name: "I64StoreConst", asm: "I64Store", argLength: 2, reg: gpstoreconst, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"},

		{name: "F32Load", asm: "F32Load", argLength: 2, reg: fpload, aux: "SymOff", typ: "Float64", symEffect: "Read"},
		{name: "F64Load", asm: "F64Load", argLength: 2, reg: fpload, aux: "SymOff", typ: "Float64", symEffect: "Read"},
		{name: "F32Store", asm: "F32Store", argLength: 3, reg: fpstore, aux: "SymOff", typ: "Mem", symEffect: "Write"},
		{name: "F64Store", asm: "F64Store", argLength: 3, reg: fpstore, aux: "SymOff", typ: "Mem", symEffect: "Write"},

		{name: "I64Const", reg: gp01, aux: "Int64", rematerializeable: true, typ: "Int64"},
		{name: "F64Const", reg: fp01, aux: "Float64", rematerializeable: true, typ: "Float64"},

		{name: "I64Eqz", asm: "I64Eqz", argLength: 1, reg: gp11, typ: "Bool"},
		{name: "I64Eq", asm: "I64Eq", argLength: 2, reg: gp21, typ: "Bool"},
		{name: "I64Ne", asm: "I64Ne", argLength: 2, reg: gp21, typ: "Bool"},
		{name: "I64LtS", asm: "I64LtS", argLength: 2, reg: gp21, typ: "Bool"},
		{name: "I64LtU", asm: "I64LtU", argLength: 2, reg: gp21, typ: "Bool"},
		{name: "I64GtS", asm: "I64GtS", argLength: 2, reg: gp21, typ: "Bool"},
		{name: "I64GtU", asm: "I64GtU", argLength: 2, reg: gp21, typ: "Bool"},
		{name: "I64LeS", asm: "I64LeS", argLength: 2, reg: gp21, typ: "Bool"},
		{name: "I64LeU", asm: "I64LeU", argLength: 2, reg: gp21, typ: "Bool"},
		{name: "I64GeS", asm: "I64GeS", argLength: 2, reg: gp21, typ: "Bool"},
		{name: "I64GeU", asm: "I64GeU", argLength: 2, reg: gp21, typ: "Bool"},

		{name: "F64Eq", asm: "F64Eq", argLength: 2, reg: fp21gp, typ: "Bool"},
		{name: "F64Ne", asm: "F64Ne", argLength: 2, reg: fp21gp, typ: "Bool"},
		{name: "F64Lt", asm: "F64Lt", argLength: 2, reg: fp21gp, typ: "Bool"},
		{name: "F64Gt", asm: "F64Gt", argLength: 2, reg: fp21gp, typ: "Bool"},
		{name: "F64Le", asm: "F64Le", argLength: 2, reg: fp21gp, typ: "Bool"},
		{name: "F64Ge", asm: "F64Ge", argLength: 2, reg: fp21gp, typ: "Bool"},

		{name: "I64Add", asm: "I64Add", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64AddConst", asm: "I64Add", argLength: 1, reg: gp11, aux: "Int64", typ: "Int64"},
		{name: "I64Sub", asm: "I64Sub", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64Mul", asm: "I64Mul", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64DivS", asm: "I64DivS", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64DivU", asm: "I64DivU", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64RemS", asm: "I64RemS", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64RemU", asm: "I64RemU", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64And", asm: "I64And", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64Or", asm: "I64Or", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64Xor", asm: "I64Xor", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64Shl", asm: "I64Shl", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64ShrS", asm: "I64ShrS", argLength: 2, reg: gp21, typ: "Int64"},
		{name: "I64ShrU", asm: "I64ShrU", argLength: 2, reg: gp21, typ: "Int64"},

		{name: "F64Neg", asm: "F64Neg", argLength: 1, reg: fp11, typ: "Float64"},
		{name: "F64Add", asm: "F64Add", argLength: 2, reg: fp21, typ: "Float64"},
		{name: "F64Sub", asm: "F64Sub", argLength: 2, reg: fp21, typ: "Float64"},
		{name: "F64Mul", asm: "F64Mul", argLength: 2, reg: fp21, typ: "Float64"},
		{name: "F64Div", asm: "F64Div", argLength: 2, reg: fp21, typ: "Float64"},

		{name: "I64TruncSF64", asm: "I64TruncSF64", argLength: 1, reg: regInfo{inputs: []regMask{fp}, outputs: []regMask{gp}}, typ: "Int64"},
		{name: "F64ConvertSI64", asm: "F64ConvertSI64", argLength: 1, reg: regInfo{inputs: []regMask{gp}, outputs: []regMask{fp}}, typ: "Float64"},
		{name: "LoweredRound32F", argLength: 1, reg: fp11, typ: "Float32"},

		{name: "Coerce", argLength: 1},
		{name: "ToUInt32", argLength: 1, typ: "UInt32"},
		{name: "ToUInt16", argLength: 1, typ: "UInt16"},
		{name: "ToUInt8", argLength: 1, typ: "UInt8"},
		{name: "ToInt32", argLength: 1, typ: "Int32"},
		{name: "ToInt16", argLength: 1, typ: "Int16"},
		{name: "ToInt8", argLength: 1, typ: "Int8"},
	}

	archs = append(archs, arch{
		name:            "Wasm",
		pkg:             "cmd/internal/obj/wasm",
		genfile:         "",
		ops:             WasmOps,
		blocks:          nil,
		regnames:        regNamesWasm,
		gpregmask:       gp,
		fpregmask:       fp,
		framepointerreg: int8(num["BP"]),
		linkreg:         -1, // not used
	})
}
