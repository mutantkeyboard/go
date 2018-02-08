// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT _rt0_wasm_js(SB),NOSPLIT,$0
    I32Const $10000 // FIXME
  Set SP

    Get SP
      Get $0 // argc
    I64ExtendUI32
  I64Store $0

    Get SP
      Get $1 // argv
    I64ExtendUI32
  I64Store $8

      I32Const runtime·rt0_go(SB)
      I32Const $16
    I32ShrU
  Set PC_F

    Get $2
  If
    Loop
          Get SP
        I32Eqz
      BrIf $1

        Get PC_F
        Get PC_B
        Get SP
      Call $0 // trace

        Get PC_F
      CallIndirect $0

      Br $0
    End
  Else
    Loop
          Get SP
        I32Eqz
      BrIf $1

        Get PC_F
      CallIndirect $0

      Br $0
    End
  End

TEXT _rt0_wasm_js_lib(SB),NOSPLIT,$0
	Unreachable
