// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT ·boolVal(SB), NOSPLIT, $0
  CallImport

TEXT ·intVal(SB), NOSPLIT, $0
  CallImport

TEXT ·floatVal(SB), NOSPLIT, $0
  CallImport

TEXT ·stringVal(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·Get(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·set(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·Index(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·setIndex(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·call(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·invoke(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·wasmnew(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·Float(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·Int(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·Bool(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·Length(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·prepareString(SB), NOSPLIT, $0
  CallImport

TEXT ·Value·loadString(SB), NOSPLIT, $0
  CallImport
