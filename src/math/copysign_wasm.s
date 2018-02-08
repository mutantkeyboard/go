// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT ·Copysign(SB),NOSPLIT,$0
    Get SP
      F64Load x+0(FP)
      F64Load y+8(FP)
    F64Copysign
	F64Store ret+16(FP)
  RET
