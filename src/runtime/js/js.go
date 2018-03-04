// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js

package js

import "unsafe"

type Value uint32

type Error struct {
	Value
}

func (err Error) Error() string {
	return "JavaScript error: " + err.Get("message").String()
}

const Undefined = Value(0)
const Null = Value(1)
const Global = Value(2)
const memory = Value(3)

var Uint8Array = Global.Get("Uint8Array")

func ValueOf(v interface{}) Value {
	switch v := v.(type) {
	case Value:
		return v
	case nil:
		return Null
	case bool:
		return boolVal(v)
	case int:
		return intVal(v)
	case int8:
		return intVal(int(v))
	case int16:
		return intVal(int(v))
	case int32:
		return intVal(int(v))
	case int64:
		return intVal(int(v))
	case uint:
		return intVal(int(v))
	case uint8:
		return intVal(int(v))
	case uint16:
		return intVal(int(v))
	case uint32:
		return intVal(int(v))
	case uint64:
		return intVal(int(v))
	case uintptr:
		return intVal(int(v))
	case unsafe.Pointer:
		return intVal(int(uintptr(v)))
	case float32:
		return floatVal(float64(v))
	case float64:
		return floatVal(v)
	case string:
		return stringVal(v)
	case []byte:
		if len(v) == 0 {
			return Uint8Array.New(memory.Get("buffer"), 0, 0)
		}
		return Uint8Array.New(memory.Get("buffer"), unsafe.Pointer(&v[0]), len(v))
	default:
		panic("invalid value")
	}
}

func boolVal(value bool) Value

func intVal(value int) Value

func floatVal(value float64) Value

func stringVal(value string) Value

func (v Value) Get(key string) Value

func (v Value) Set(key string, value interface{}) {
	v.set(key, ValueOf(value))
}

func (v Value) set(key string, value Value)

func (v Value) Index(i int) Value

func (v Value) SetIndex(i int, value interface{}) {
	v.setIndex(i, ValueOf(value))
}

func (v Value) setIndex(i int, value Value)

func makeArgs(args []interface{}) []Value {
	argVals := make([]Value, len(args))
	for i, arg := range args {
		argVals[i] = ValueOf(arg)
	}
	return argVals
}

func (v Value) Call(name string, args ...interface{}) Value {
	res, ok := v.call(name, makeArgs(args))
	if !ok {
		panic(Error{res})
	}
	return res
}

func (v Value) call(name string, args []Value) (Value, bool)

func (v Value) Invoke(args ...interface{}) Value {
	res, ok := v.invoke(makeArgs(args))
	if !ok {
		panic(Error{res})
	}
	return res
}

func (v Value) invoke(args []Value) (Value, bool)

func (v Value) New(args ...interface{}) Value {
	res, ok := v.wasmnew(makeArgs(args))
	if !ok {
		panic(Error{res})
	}
	return res
}

func (v Value) wasmnew(args []Value) (Value, bool)

func (v Value) Float() float64

func (v Value) Int() int

func (v Value) Bool() bool

func (v Value) Length() int

func (v Value) String() string {
	str, length := v.prepareString()
	b := make([]byte, length)
	str.loadString(b)
	return string(b)
}

func (v Value) prepareString() (Value, int)

func (v Value) loadString(b []byte)
