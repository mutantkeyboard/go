// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT runtime·wasmcall(SB), NOSPLIT, $0-0
      Get SP
      I32Const $8
    I32Sub
	Set SP

	  Get SP
	    Get $1
    I64ExtendUI32
	I64Store $0

	  Get SP
    Get PC_F
	I32Store16 $2

	  Get $0
	Set PC_F

	  I32Const $0
	Set PC_B

TEXT runtime·wasmreturn(SB), NOSPLIT, $0-0
	    Get SP
	  I32Load16U $0
	Set PC_B

      Get SP
    I32Load16U $2
	Set PC_F

	    Get SP
	    I32Const $8
	  I32Add
	Set SP

TEXT runtime·wasmmove(SB), NOSPLIT, $0-0
  Block
    Loop
      // *dst = *src
        Get $0
          Get $1
        I64Load $0
      I64Store $0

      // n--
          Get $2
          I32Const $1
        I32Sub
      Set $2

      // n == 0
          Get $2
        I32Eqz
      BrIf $1

      // dst += 8
          Get $0
          I32Const $8
        I32Add
      Set $0

      // src += 8
          Get $1
          I32Const $8
        I32Add
      Set $1

      Br $0
    End
  End

TEXT runtime·wasmzero(SB), NOSPLIT, $0-0
  Block
    Loop
      // *dst = 0
        Get $0
        I64Const $0
      I64Store $0

      // n--
          Get $1
          I32Const $1
        I32Sub
      Set $1

      // n == 0
          Get $1
        I32Eqz
      BrIf $1

      // dst += 8
          Get $0
          I32Const $8
        I32Add
      Set $0

      Br $0
    End
  End

TEXT runtime·wasmdiv(SB), NOSPLIT, $0-0
      Get $0
      I64Const $-0x8000000000000000
    I64Eq
  If
        Get $1
        I64Const $-1
      I64Eq
    If
      I64Const $-0x8000000000000000
      Br $2
    End
  End
  Get $0
  Get $1
  I64DivS

TEXT runtime·wasmtrunc(SB), NOSPLIT, $0-0
      Get $0
      Get $0
    F64Ne // NaN
  If
    I64Const $0x8000000000000000
    Br $1
  End
      Get $0
      F64Const $9223372036854775807.
    F64Gt
  If
    I64Const $0x8000000000000000
    Br $1
  End
      Get $0
      F64Const $-9223372036854775808.
    F64Lt
  If
    I64Const $0x8000000000000000
    Br $1
  End
    Get $0
  I64TruncSF64

TEXT runtime·exit(SB), NOSPLIT, $0-8
  Call runtime·wasmexit(SB)
    I32Const $0
  Set SP

TEXT runtime·exitThread(SB), NOSPLIT, $0-0
  Unreachable

TEXT runtime·getclosureptr(SB), NOSPLIT, $0-0
  Unreachable

TEXT runtime·osyield(SB), NOSPLIT, $0-0
  Unreachable

TEXT runtime·usleep(SB), NOSPLIT, $0-0
  RET // FIXME

TEXT runtime·IndexByte(SB), NOSPLIT, $0-0
  Unreachable

TEXT runtime·growMemory(SB), NOSPLIT, $0-12
    Get SP
        I32Load pages+0(FP)
        CurrentMemory
      I32Sub
    GrowMemory
  I32Store ret+8(FP)
  RET

TEXT ·wasmexit(SB), NOSPLIT, $0
  CallImport

TEXT ·wasmwrite(SB), NOSPLIT, $0
  CallImport

TEXT ·nanotime(SB), NOSPLIT, $0
  CallImport

TEXT ·walltime(SB), NOSPLIT, $0
  CallImport
