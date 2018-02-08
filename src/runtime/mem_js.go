// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"unsafe"
)

var allocEnd uintptr

// Don't split the stack as this function may be invoked without a valid G,
// which prevents us from allocating more stack.
//go:nosplit
func sysAlloc(n uintptr, sysStat *uint64) unsafe.Pointer {
	// println("sysAlloc", n, uintptr(allocEnd))
	p := unsafe.Pointer(allocEnd)
	allocEnd += n
	mSysStatInc(sysStat, n)
	return p
}

func sysUnused(v unsafe.Pointer, n uintptr) {
}

func sysUsed(v unsafe.Pointer, n uintptr) {
}

// Don't split the stack as this function may be invoked without a valid G,
// which prevents us from allocating more stack.
//go:nosplit
func sysFree(v unsafe.Pointer, n uintptr, sysStat *uint64) {
}

func sysFault(v unsafe.Pointer, n uintptr) {
}

func sysReserve(v unsafe.Pointer, n uintptr, reserved *bool) unsafe.Pointer {
	if n > 1024*1024*1024 {
		return nil
	}
	// println("sysReserve", uintptr(v), n)
	// TODO grow_memory
	allocEnd = uintptr(v) + n
	// growMemory(int32(allocEnd / sys.DefaultPhysPageSize))
	*reserved = true
	return v
}

func growMemory(pages int32) int32

func sysMap(v unsafe.Pointer, n uintptr, reserved bool, sysStat *uint64) {
	mSysStatInc(sysStat, n)
}
