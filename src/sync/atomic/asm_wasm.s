// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT ·SwapInt32(SB),NOSPLIT,$0-20
	JMP	·SwapUint32(SB)

TEXT ·SwapUint32(SB),NOSPLIT,$0-20
	  I64Load addr+0(FP)
  Set I0

    Get SP
    I32Load (I0)
  I32Store old+16(FP)

      Get I0
    I32WrapI64
    I32Load new+8(FP)
  I32Store $0

  RET

TEXT ·SwapInt64(SB),NOSPLIT,$0-24
	JMP	·SwapUint64(SB)

TEXT ·SwapUint64(SB),NOSPLIT,$0-24
	  I64Load addr+0(FP)
  Set I0

    Get SP
    I64Load (I0)
  I64Store old+16(FP)

      Get I0
    I32WrapI64
    I64Load new+8(FP)
  I64Store $0

  RET

TEXT ·SwapUintptr(SB),NOSPLIT,$0-24
	JMP	·SwapUint64(SB)

TEXT ·CompareAndSwapInt32(SB),NOSPLIT,$0-17
	JMP	·CompareAndSwapUint32(SB)

TEXT ·CompareAndSwapUint32(SB),NOSPLIT,$0-17
  Block
      I64Load addr+0(FP)
    Set I0

        I32Load (I0)
        I32Load old+8(FP)
      I32Ne
    If
              Get SP
        I64Const $0
      I64Store8 swapped+16(FP)
      Br $1
    End

        Get I0
      I32WrapI64
      I32Load new+12(FP)
    I32Store $0

      Get SP
      I64Const $1
    I64Store8 swapped+16(FP)
  End

  RET

TEXT ·CompareAndSwapUintptr(SB),NOSPLIT,$0-25
	JMP	·CompareAndSwapUint64(SB)

TEXT ·CompareAndSwapInt64(SB),NOSPLIT,$0-25
	JMP	·CompareAndSwapUint64(SB)

TEXT ·CompareAndSwapUint64(SB),NOSPLIT,$0-25
  Block
      I64Load addr+0(FP)
    Set I0

        I64Load (I0)
        I64Load old+8(FP)
      I64Ne
    If
              Get SP
        I64Const $0
      I64Store8 swapped+24(FP)
      Br $1
    End

        Get I0
      I32WrapI64
      I64Load new+16(FP)
    I64Store $0

      Get SP
      I64Const $1
    I64Store8 swapped+24(FP)
  End

  RET

TEXT ·AddInt32(SB),NOSPLIT,$0-20
	JMP	·AddUint32(SB)

TEXT ·AddUint32(SB),NOSPLIT,$0-20
    I64Load addr+0(FP)
  Set I0

      I64Load32U (I0)
      I64Load32U delta+8(FP)
    I64Add
  Set I1

      Get I0
    I32WrapI64
    Get I1
  I64Store32 $0

    Get SP
    Get I1
  I64Store32 new+16(FP)

	RET

TEXT ·AddUintptr(SB),NOSPLIT,$0-24
	JMP	·AddUint64(SB)

TEXT ·AddInt64(SB),NOSPLIT,$0-24
	JMP	·AddUint64(SB)

TEXT ·AddUint64(SB),NOSPLIT,$0-24
    I64Load addr+0(FP)
  Set I0

      I64Load (I0)
      I64Load delta+8(FP)
    I64Add
  Set I1

      Get I0
    I32WrapI64
    Get I1
  I64Store $0

    Get SP
    Get I1
  I64Store new+16(FP)

	RET

TEXT ·LoadInt32(SB),NOSPLIT,$0-12
	JMP	·LoadUint32(SB)

TEXT ·LoadUint32(SB),NOSPLIT,$0-12
    Get SP
        I64Load addr+0(FP)
      I32WrapI64
    I64Load32U $0
  I64Store32 val+8(FP)
	RET

TEXT ·LoadInt64(SB),NOSPLIT,$0-16
	JMP	·LoadUint64(SB)

TEXT ·LoadUintptr(SB),NOSPLIT,$0-16
	JMP	·LoadUint64(SB)

TEXT ·LoadPointer(SB),NOSPLIT,$0-16
  JMP	·LoadUint64(SB)

TEXT ·LoadUint64(SB),NOSPLIT,$0-16
    Get SP
        I64Load addr+0(FP)
      I32WrapI64
    I64Load $0
  I64Store val+8(FP)
	RET

TEXT ·StoreInt32(SB),NOSPLIT,$0-12
	JMP	·StoreUint32(SB)

TEXT ·StoreUint32(SB),NOSPLIT,$0-12
      I64Load addr+0(FP)
    I32WrapI64
    I32Load val+8(FP)
  I32Store $0
  RET

TEXT ·StoreInt64(SB),NOSPLIT,$0-16
	JMP	·StoreUint64(SB)

TEXT ·StoreUint64(SB),NOSPLIT,$0-16
      I64Load addr+0(FP)
    I32WrapI64
    I64Load val+8(FP)
  I64Store $0
  RET

TEXT ·StoreUintptr(SB),NOSPLIT,$0-16
	JMP	·StoreUint64(SB)
