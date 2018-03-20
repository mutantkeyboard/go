// Code generated by mkbuiltin.go. DO NOT EDIT.

package gc

import "cmd/compile/internal/types"

var runtimeDecls = [...]struct {
	name string
	tag  int
	typ  int
}{
	{"newobject", funcTag, 4},
	{"panicindex", funcTag, 5},
	{"panicslice", funcTag, 5},
	{"panicdivide", funcTag, 5},
	{"throwinit", funcTag, 5},
	{"panicwrap", funcTag, 5},
	{"gopanic", funcTag, 7},
	{"gorecover", funcTag, 10},
	{"goschedguarded", funcTag, 5},
	{"printbool", funcTag, 12},
	{"printfloat", funcTag, 14},
	{"printint", funcTag, 16},
	{"printhex", funcTag, 18},
	{"printuint", funcTag, 18},
	{"printcomplex", funcTag, 20},
	{"printstring", funcTag, 22},
	{"printpointer", funcTag, 23},
	{"printiface", funcTag, 23},
	{"printeface", funcTag, 23},
	{"printslice", funcTag, 23},
	{"printnl", funcTag, 5},
	{"printsp", funcTag, 5},
	{"printlock", funcTag, 5},
	{"printunlock", funcTag, 5},
	{"concatstring2", funcTag, 26},
	{"concatstring3", funcTag, 27},
	{"concatstring4", funcTag, 28},
	{"concatstring5", funcTag, 29},
	{"concatstrings", funcTag, 31},
	{"cmpstring", funcTag, 33},
	{"intstring", funcTag, 36},
	{"slicebytetostring", funcTag, 38},
	{"slicebytetostringtmp", funcTag, 39},
	{"slicerunetostring", funcTag, 42},
	{"stringtoslicebyte", funcTag, 43},
	{"stringtoslicerune", funcTag, 46},
	{"decoderune", funcTag, 47},
	{"slicecopy", funcTag, 49},
	{"slicestringcopy", funcTag, 50},
	{"convI2I", funcTag, 51},
	{"convT2E", funcTag, 52},
	{"convT2E16", funcTag, 52},
	{"convT2E32", funcTag, 52},
	{"convT2E64", funcTag, 52},
	{"convT2Estring", funcTag, 52},
	{"convT2Eslice", funcTag, 52},
	{"convT2Enoptr", funcTag, 52},
	{"convT2I", funcTag, 52},
	{"convT2I16", funcTag, 52},
	{"convT2I32", funcTag, 52},
	{"convT2I64", funcTag, 52},
	{"convT2Istring", funcTag, 52},
	{"convT2Islice", funcTag, 52},
	{"convT2Inoptr", funcTag, 52},
	{"assertE2I", funcTag, 51},
	{"assertE2I2", funcTag, 53},
	{"assertI2I", funcTag, 51},
	{"assertI2I2", funcTag, 53},
	{"panicdottypeE", funcTag, 54},
	{"panicdottypeI", funcTag, 54},
	{"panicnildottype", funcTag, 55},
	{"ifaceeq", funcTag, 58},
	{"efaceeq", funcTag, 58},
	{"fastrand", funcTag, 60},
	{"makemap64", funcTag, 62},
	{"makemap", funcTag, 63},
	{"makemap_small", funcTag, 64},
	{"mapaccess1", funcTag, 65},
	{"mapaccess1_fast32", funcTag, 66},
	{"mapaccess1_fast64", funcTag, 66},
	{"mapaccess1_faststr", funcTag, 66},
	{"mapaccess1_fat", funcTag, 67},
	{"mapaccess2", funcTag, 68},
	{"mapaccess2_fast32", funcTag, 69},
	{"mapaccess2_fast64", funcTag, 69},
	{"mapaccess2_faststr", funcTag, 69},
	{"mapaccess2_fat", funcTag, 70},
	{"mapassign", funcTag, 65},
	{"mapassign_fast32", funcTag, 66},
	{"mapassign_fast32ptr", funcTag, 66},
	{"mapassign_fast64", funcTag, 66},
	{"mapassign_fast64ptr", funcTag, 66},
	{"mapassign_faststr", funcTag, 66},
	{"mapiterinit", funcTag, 71},
	{"mapdelete", funcTag, 71},
	{"mapdelete_fast32", funcTag, 72},
	{"mapdelete_fast64", funcTag, 72},
	{"mapdelete_faststr", funcTag, 72},
	{"mapiternext", funcTag, 73},
	{"makechan64", funcTag, 75},
	{"makechan", funcTag, 76},
	{"chanrecv1", funcTag, 78},
	{"chanrecv2", funcTag, 79},
	{"chansend1", funcTag, 81},
	{"closechan", funcTag, 23},
	{"writeBarrier", varTag, 83},
	{"typedmemmove", funcTag, 84},
	{"typedmemclr", funcTag, 85},
	{"typedslicecopy", funcTag, 86},
	{"selectnbsend", funcTag, 87},
	{"selectnbrecv", funcTag, 88},
	{"selectnbrecv2", funcTag, 90},
	{"newselect", funcTag, 91},
	{"selectsend", funcTag, 92},
	{"selectrecv", funcTag, 93},
	{"selectdefault", funcTag, 55},
	{"selectgo", funcTag, 94},
	{"block", funcTag, 5},
	{"makeslice", funcTag, 96},
	{"makeslice64", funcTag, 97},
	{"growslice", funcTag, 98},
	{"memmove", funcTag, 99},
	{"memclrNoHeapPointers", funcTag, 100},
	{"memclrHasPointers", funcTag, 100},
	{"memequal", funcTag, 101},
	{"memequal8", funcTag, 102},
	{"memequal16", funcTag, 102},
	{"memequal32", funcTag, 102},
	{"memequal64", funcTag, 102},
	{"memequal128", funcTag, 102},
	{"int64div", funcTag, 103},
	{"uint64div", funcTag, 104},
	{"int64mod", funcTag, 103},
	{"uint64mod", funcTag, 104},
	{"float64toint64", funcTag, 105},
	{"float64touint64", funcTag, 106},
	{"float64touint32", funcTag, 107},
	{"int64tofloat64", funcTag, 108},
	{"uint64tofloat64", funcTag, 109},
	{"uint32tofloat64", funcTag, 110},
	{"complex128div", funcTag, 111},
	{"racefuncenter", funcTag, 112},
	{"racefuncexit", funcTag, 5},
	{"raceread", funcTag, 112},
	{"racewrite", funcTag, 112},
	{"racereadrange", funcTag, 113},
	{"racewriterange", funcTag, 113},
	{"msanread", funcTag, 113},
	{"msanwrite", funcTag, 113},
	{"support_popcnt", varTag, 11},
	{"support_sse41", varTag, 11},
}

func runtimeTypes() []*types.Type {
	var typs [114]*types.Type
	typs[0] = types.Bytetype
	typs[1] = types.NewPtr(typs[0])
	typs[2] = types.Types[TANY]
	typs[3] = types.NewPtr(typs[2])
	typs[4] = functype(nil, []*Node{anonfield(typs[1])}, []*Node{anonfield(typs[3])})
	typs[5] = functype(nil, nil, nil)
	typs[6] = types.Types[TINTER]
	typs[7] = functype(nil, []*Node{anonfield(typs[6])}, nil)
	typs[8] = types.Types[TINT32]
	typs[9] = types.NewPtr(typs[8])
	typs[10] = functype(nil, []*Node{anonfield(typs[9])}, []*Node{anonfield(typs[6])})
	typs[11] = types.Types[TBOOL]
	typs[12] = functype(nil, []*Node{anonfield(typs[11])}, nil)
	typs[13] = types.Types[TFLOAT64]
	typs[14] = functype(nil, []*Node{anonfield(typs[13])}, nil)
	typs[15] = types.Types[TINT64]
	typs[16] = functype(nil, []*Node{anonfield(typs[15])}, nil)
	typs[17] = types.Types[TUINT64]
	typs[18] = functype(nil, []*Node{anonfield(typs[17])}, nil)
	typs[19] = types.Types[TCOMPLEX128]
	typs[20] = functype(nil, []*Node{anonfield(typs[19])}, nil)
	typs[21] = types.Types[TSTRING]
	typs[22] = functype(nil, []*Node{anonfield(typs[21])}, nil)
	typs[23] = functype(nil, []*Node{anonfield(typs[2])}, nil)
	typs[24] = types.NewArray(typs[0], 32)
	typs[25] = types.NewPtr(typs[24])
	typs[26] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[27] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[28] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[29] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[30] = types.NewSlice(typs[21])
	typs[31] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[30])}, []*Node{anonfield(typs[21])})
	typs[32] = types.Types[TINT]
	typs[33] = functype(nil, []*Node{anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[32])})
	typs[34] = types.NewArray(typs[0], 4)
	typs[35] = types.NewPtr(typs[34])
	typs[36] = functype(nil, []*Node{anonfield(typs[35]), anonfield(typs[15])}, []*Node{anonfield(typs[21])})
	typs[37] = types.NewSlice(typs[0])
	typs[38] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[37])}, []*Node{anonfield(typs[21])})
	typs[39] = functype(nil, []*Node{anonfield(typs[37])}, []*Node{anonfield(typs[21])})
	typs[40] = types.Runetype
	typs[41] = types.NewSlice(typs[40])
	typs[42] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[41])}, []*Node{anonfield(typs[21])})
	typs[43] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21])}, []*Node{anonfield(typs[37])})
	typs[44] = types.NewArray(typs[40], 32)
	typs[45] = types.NewPtr(typs[44])
	typs[46] = functype(nil, []*Node{anonfield(typs[45]), anonfield(typs[21])}, []*Node{anonfield(typs[41])})
	typs[47] = functype(nil, []*Node{anonfield(typs[21]), anonfield(typs[32])}, []*Node{anonfield(typs[40]), anonfield(typs[32])})
	typs[48] = types.Types[TUINTPTR]
	typs[49] = functype(nil, []*Node{anonfield(typs[2]), anonfield(typs[2]), anonfield(typs[48])}, []*Node{anonfield(typs[32])})
	typs[50] = functype(nil, []*Node{anonfield(typs[2]), anonfield(typs[2])}, []*Node{anonfield(typs[32])})
	typs[51] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[2])}, []*Node{anonfield(typs[2])})
	typs[52] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3])}, []*Node{anonfield(typs[2])})
	typs[53] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[2])}, []*Node{anonfield(typs[2]), anonfield(typs[11])})
	typs[54] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[1]), anonfield(typs[1])}, nil)
	typs[55] = functype(nil, []*Node{anonfield(typs[1])}, nil)
	typs[56] = types.NewPtr(typs[48])
	typs[57] = types.Types[TUNSAFEPTR]
	typs[58] = functype(nil, []*Node{anonfield(typs[56]), anonfield(typs[57]), anonfield(typs[57])}, []*Node{anonfield(typs[11])})
	typs[59] = types.Types[TUINT32]
	typs[60] = functype(nil, nil, []*Node{anonfield(typs[59])})
	typs[61] = types.NewMap(typs[2], typs[2])
	typs[62] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15]), anonfield(typs[3])}, []*Node{anonfield(typs[61])})
	typs[63] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[32]), anonfield(typs[3])}, []*Node{anonfield(typs[61])})
	typs[64] = functype(nil, nil, []*Node{anonfield(typs[61])})
	typs[65] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[61]), anonfield(typs[3])}, []*Node{anonfield(typs[3])})
	typs[66] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[61]), anonfield(typs[2])}, []*Node{anonfield(typs[3])})
	typs[67] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[61]), anonfield(typs[3]), anonfield(typs[1])}, []*Node{anonfield(typs[3])})
	typs[68] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[61]), anonfield(typs[3])}, []*Node{anonfield(typs[3]), anonfield(typs[11])})
	typs[69] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[61]), anonfield(typs[2])}, []*Node{anonfield(typs[3]), anonfield(typs[11])})
	typs[70] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[61]), anonfield(typs[3]), anonfield(typs[1])}, []*Node{anonfield(typs[3]), anonfield(typs[11])})
	typs[71] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[61]), anonfield(typs[3])}, nil)
	typs[72] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[61]), anonfield(typs[2])}, nil)
	typs[73] = functype(nil, []*Node{anonfield(typs[3])}, nil)
	typs[74] = types.NewChan(typs[2], types.Cboth)
	typs[75] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15])}, []*Node{anonfield(typs[74])})
	typs[76] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[32])}, []*Node{anonfield(typs[74])})
	typs[77] = types.NewChan(typs[2], types.Crecv)
	typs[78] = functype(nil, []*Node{anonfield(typs[77]), anonfield(typs[3])}, nil)
	typs[79] = functype(nil, []*Node{anonfield(typs[77]), anonfield(typs[3])}, []*Node{anonfield(typs[11])})
	typs[80] = types.NewChan(typs[2], types.Csend)
	typs[81] = functype(nil, []*Node{anonfield(typs[80]), anonfield(typs[3])}, nil)
	typs[82] = types.NewArray(typs[0], 3)
	typs[83] = tostruct([]*Node{namedfield("enabled", typs[11]), namedfield("pad", typs[82]), namedfield("needed", typs[11]), namedfield("cgo", typs[11]), namedfield("alignme", typs[17])})
	typs[84] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3]), anonfield(typs[3])}, nil)
	typs[85] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3])}, nil)
	typs[86] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[2]), anonfield(typs[2])}, []*Node{anonfield(typs[32])})
	typs[87] = functype(nil, []*Node{anonfield(typs[80]), anonfield(typs[3])}, []*Node{anonfield(typs[11])})
	typs[88] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[77])}, []*Node{anonfield(typs[11])})
	typs[89] = types.NewPtr(typs[11])
	typs[90] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[89]), anonfield(typs[77])}, []*Node{anonfield(typs[11])})
	typs[91] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15]), anonfield(typs[8])}, nil)
	typs[92] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[80]), anonfield(typs[3])}, nil)
	typs[93] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[77]), anonfield(typs[3]), anonfield(typs[89])}, nil)
	typs[94] = functype(nil, []*Node{anonfield(typs[1])}, []*Node{anonfield(typs[32])})
	typs[95] = types.NewSlice(typs[2])
	typs[96] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[32]), anonfield(typs[32])}, []*Node{anonfield(typs[95])})
	typs[97] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15]), anonfield(typs[15])}, []*Node{anonfield(typs[95])})
	typs[98] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[95]), anonfield(typs[32])}, []*Node{anonfield(typs[95])})
	typs[99] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[3]), anonfield(typs[48])}, nil)
	typs[100] = functype(nil, []*Node{anonfield(typs[57]), anonfield(typs[48])}, nil)
	typs[101] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[3]), anonfield(typs[48])}, []*Node{anonfield(typs[11])})
	typs[102] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[3])}, []*Node{anonfield(typs[11])})
	typs[103] = functype(nil, []*Node{anonfield(typs[15]), anonfield(typs[15])}, []*Node{anonfield(typs[15])})
	typs[104] = functype(nil, []*Node{anonfield(typs[17]), anonfield(typs[17])}, []*Node{anonfield(typs[17])})
	typs[105] = functype(nil, []*Node{anonfield(typs[13])}, []*Node{anonfield(typs[15])})
	typs[106] = functype(nil, []*Node{anonfield(typs[13])}, []*Node{anonfield(typs[17])})
	typs[107] = functype(nil, []*Node{anonfield(typs[13])}, []*Node{anonfield(typs[59])})
	typs[108] = functype(nil, []*Node{anonfield(typs[15])}, []*Node{anonfield(typs[13])})
	typs[109] = functype(nil, []*Node{anonfield(typs[17])}, []*Node{anonfield(typs[13])})
	typs[110] = functype(nil, []*Node{anonfield(typs[59])}, []*Node{anonfield(typs[13])})
	typs[111] = functype(nil, []*Node{anonfield(typs[19]), anonfield(typs[19])}, []*Node{anonfield(typs[19])})
	typs[112] = functype(nil, []*Node{anonfield(typs[48])}, nil)
	typs[113] = functype(nil, []*Node{anonfield(typs[48]), anonfield(typs[48])}, nil)
	return typs[:]
}
