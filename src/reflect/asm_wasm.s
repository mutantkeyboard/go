// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "funcdata.h"

// makeFuncStub is the code half of the function returned by MakeFunc.
// See the comment on the declaration of makeFuncStub in makefunc.go
// for more details.
// No arg size here; runtime pulls arg map out of the func value.
TEXT ·makeFuncStub(SB),(NOSPLIT|WRAPPER),$16
      NO_LOCAL_POINTERS
        Get SP
        Get CTX
      I64Store $0
        Get SP
            Get SP
          I64ExtendUI32
          I64Const $argframe+0(FP)
        I64Add
      I64Store $8
      CALL ·callReflect(SB)
    End
  End

// methodValueCall is the code half of the function returned by makeMethodValue.
// See the comment on the declaration of methodValueCall in makefunc.go
// for more details.
// No arg size here; runtime pulls arg map out of the func value.
TEXT ·methodValueCall(SB),(NOSPLIT|WRAPPER),$16
      NO_LOCAL_POINTERS
        Get SP
        Get CTX
      I64Store $0
        Get SP
            Get SP
          I64ExtendUI32
          I64Const $argframe+0(FP)
        I64Add
      I64Store $8
      CALL ·callMethod(SB)
    End
  End
