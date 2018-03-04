// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

import "runtime/js"

func init() {
	Reader = &reader{}
}

type reader struct{}

var jsCrypto = js.Global.Get("crypto")

func (r *reader) Read(b []byte) (int, error) {
	jsCrypto.Call("getRandomValues", b)
	return len(b), nil
}
