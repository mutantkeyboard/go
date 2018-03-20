// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js

package js_test

import (
	"runtime/js"
	"testing"
)

var dummys = js.Global.Call("eval", `({
	someBool: true,
	someString: "abc\u1234",
	someInt: 42,
	someFloat: 42.123,
	someArray: [41, 42, 43],
	add: function(a, b) {
		return a + b;
	},
})`)

func TestBool(t *testing.T) {
	e := true
	o := dummys.Get("someBool")
	if v := o.Bool(); v != e {
		t.Errorf("expected %#v, got %#v", e, v)
	}
	if dummys.Set("otherBool", e); dummys.Get("otherBool").Bool() != e {
		t.Fail()
	}
}

func TestString(t *testing.T) {
	e := "abc\u1234"
	o := dummys.Get("someString")
	if v := o.String(); v != e {
		t.Errorf("expected %#v, got %#v", e, v)
	}
	if dummys.Set("otherString", e); dummys.Get("otherString").String() != e {
		t.Fail()
	}
}

func TestInt(t *testing.T) {
	e := 42
	o := dummys.Get("someInt")
	if v := o.Int(); v != e {
		t.Errorf("expected %#v, got %#v", e, v)
	}
	if dummys.Set("otherInt", e); dummys.Get("otherInt").Int() != e {
		t.Fail()
	}
}

func TestFloat(t *testing.T) {
	e := 42.123
	o := dummys.Get("someFloat")
	if v := o.Float(); v != e {
		t.Errorf("expected %#v, got %#v", e, v)
	}
	if dummys.Set("otherFloat", e); dummys.Get("otherFloat").Float() != e {
		t.Fail()
	}
}

func TestUndefined(t *testing.T) {
	dummys.Set("test", js.Undefined)
	if dummys == js.Undefined || dummys.Get("test") != js.Undefined || dummys.Get("xyz") != js.Undefined {
		t.Fail()
	}
}

func TestNull(t *testing.T) {
	dummys.Set("test1", nil)
	dummys.Set("test2", js.Null)
	if dummys == js.Null || dummys.Get("test1") != js.Null || dummys.Get("test2") != js.Null {
		t.Fail()
	}
}

func TestLength(t *testing.T) {
	if dummys.Get("someArray").Length() != 3 {
		t.Fail()
	}
}

func TestIndex(t *testing.T) {
	if dummys.Get("someArray").Index(1).Int() != 42 {
		t.Fail()
	}
}

func TestSetIndex(t *testing.T) {
	dummys.Get("someArray").SetIndex(2, 99)
	if dummys.Get("someArray").Index(2).Int() != 99 {
		t.Fail()
	}
}

func TestCall(t *testing.T) {
	var i int64 = 40
	if dummys.Call("add", i, 2).Int() != 42 {
		t.Fail()
	}
	if dummys.Call("add", js.Global.Call("eval", "40"), 2).Int() != 42 {
		t.Fail()
	}
}

func TestInvoke(t *testing.T) {
	var i int64 = 40
	if dummys.Get("add").Invoke(i, 2).Int() != 42 {
		t.Fail()
	}
}

func TestNew(t *testing.T) {
	if js.Global.Get("Array").New(42).Length() != 42 {
		t.Fail()
	}
}
<<<<<<< HEAD

// func TestFunc(t *testing.T) {
// 	a := dummys.Call("mapArray", []int{1, 2, 3}, func(e int64) int64 { return e + 40 })
// 	b := dummys.Call("mapArray", []int{1, 2, 3}, func(e ...int64) int64 { return e[0] + 40 })
// 	if a.Index(1).Int() != 42 || b.Index(1).Int() != 42 {
// 		t.Fail()
// 	}

// 	add := dummys.Get("add").Interface().(func(...interface{}) *js.Object)
// 	var i int64 = 40
// 	if add(i, 2).Int() != 42 {
// 		t.Fail()
// 	}
// }

// func TestDate(t *testing.T) {
// 	d := time.Date(2013, time.August, 27, 22, 25, 11, 0, time.UTC)
// 	if dummys.Call("toUnixTimestamp", d).Int() != int(d.Unix()) {
// 		t.Fail()
// 	}

// 	d2 := js.Global.Get("Date").New(d.UnixNano() / 1000000).Interface().(time.Time)
// 	if !d2.Equal(d) {
// 		t.Fail()
// 	}
// }

// // https://github.com/gopherjs/gopherjs/issues/287
// func TestInternalizeDate(t *testing.T) {
// 	var a = time.Unix(0, (123 * time.Millisecond).Nanoseconds())
// 	var b time.Time
// 	js.Global.Set("internalizeDate", func(t time.Time) { b = t })
// 	js.Global.Call("eval", "(internalizeDate(new Date(123)))")
// 	if a != b {
// 		t.Fail()
// 	}
// }

// func TestEquality(t *testing.T) {
// 	if js.Global.Get("Array") != js.Global.Get("Array") || js.Global.Get("Array") == js.Global.Get("String") {
// 		t.Fail()
// 	}
// 	type S struct{ *js.Object }
// 	o1 := js.Global.Get("Object").New()
// 	o2 := js.Global.Get("Object").New()
// 	a := S{o1}
// 	b := S{o1}
// 	c := S{o2}
// 	if a != b || a == c {
// 		t.Fail()
// 	}
// }

// func TestUndefinedEquality(t *testing.T) {
// 	var ui interface{} = js.Undefined
// 	if ui != js.Undefined {
// 		t.Fail()
// 	}
// }

// func TestInterfaceEquality(t *testing.T) {
// 	o := js.Global.Get("Object").New()
// 	var i interface{} = o
// 	if i != o {
// 		t.Fail()
// 	}
// }

// func TestUndefinedInternalization(t *testing.T) {
// 	undefinedEqualsJsUndefined := func(i interface{}) bool {
// 		return i == js.Undefined
// 	}
// 	js.Global.Set("undefinedEqualsJsUndefined", undefinedEqualsJsUndefined)
// 	if !js.Global.Call("eval", "(undefinedEqualsJsUndefined(undefined))").Bool() {
// 		t.Fail()
// 	}
// }

// func TestSameFuncWrapper(t *testing.T) {
// 	a := func(_ string) {} // string argument to force wrapping
// 	b := func(_ string) {} // string argument to force wrapping
// 	if !dummys.Call("isEqual", a, a).Bool() || dummys.Call("isEqual", a, b).Bool() {
// 		t.Fail()
// 	}
// 	if !dummys.Call("isEqual", somePackageFunction, somePackageFunction).Bool() {
// 		t.Fail()
// 	}
// 	if !dummys.Call("isEqual", (*T).someMethod, (*T).someMethod).Bool() {
// 		t.Fail()
// 	}
// 	t1 := &T{}
// 	t2 := &T{}
// 	if !dummys.Call("isEqual", t1.someMethod, t1.someMethod).Bool() || dummys.Call("isEqual", t1.someMethod, t2.someMethod).Bool() {
// 		t.Fail()
// 	}
// }

// func somePackageFunction(_ string) {
// }

// type T struct{}

// func (t *T) someMethod() {
// 	println(42)
// }

// func TestError(t *testing.T) {
// 	defer func() {
// 		err := recover()
// 		if err == nil {
// 			t.Fail()
// 		}
// 		if _, ok := err.(error); !ok {
// 			t.Fail()
// 		}
// 		jsErr, ok := err.(*js.Error)
// 		if !ok || !strings.Contains(jsErr.Error(), "throwsError") {
// 			t.Fail()
// 		}
// 	}()
// 	js.Global.Get("notExisting").Call("throwsError")
// }

// type F struct {
// 	Field int
// }

// func TestExternalizeField(t *testing.T) {
// 	if dummys.Call("testField", map[string]int{"Field": 42}).Int() != 42 {
// 		t.Fail()
// 	}
// 	if dummys.Call("testField", F{42}).Int() != 42 {
// 		t.Fail()
// 	}
// }

// func TestMakeFunc(t *testing.T) {
// 	o := js.Global.Get("Object").New()
// 	for i := 3; i < 5; i++ {
// 		x := i
// 		if i == 4 {
// 			break
// 		}
// 		o.Set("f", js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
// 			if this != o {
// 				t.Fail()
// 			}
// 			if len(arguments) != 2 || arguments[0].Int() != 1 || arguments[1].Int() != 2 {
// 				t.Fail()
// 			}
// 			return x
// 		}))
// 	}
// 	if o.Call("f", 1, 2).Int() != 3 {
// 		t.Fail()
// 	}
// }

// type M struct {
// 	f int
// }

// func (m *M) Method(a interface{}) map[string]string {
// 	if a.(map[string]interface{})["x"].(float64) != 1 || m.f != 42 {
// 		return nil
// 	}
// 	return map[string]string{
// 		"y": "z",
// 	}
// }

// func TestMakeWrapper(t *testing.T) {
// 	m := &M{42}
// 	if !js.Global.Call("eval", `(function(m) { return m.Method({x: 1})["y"] === "z"; })`).Invoke(js.MakeWrapper(m)).Bool() {
// 		t.Fail()
// 	}

// 	if js.MakeWrapper(m).Interface() != m {
// 		t.Fail()
// 	}

// 	f := func(m *M) {
// 		if m.f != 42 {
// 			t.Fail()
// 		}
// 	}
// 	js.Global.Call("eval", `(function(f, m) { f(m); })`).Invoke(f, js.MakeWrapper(m))
// }

// func TestCallWithNull(t *testing.T) {
// 	c := make(chan int, 1)
// 	js.Global.Set("test", func() {
// 		c <- 42
// 	})
// 	js.Global.Get("test").Call("call", nil)
// 	if <-c != 42 {
// 		t.Fail()
// 	}
// }

// func TestReflection(t *testing.T) {
// 	o := js.Global.Call("eval", "({ answer: 42 })")
// 	if reflect.ValueOf(o).Interface().(*js.Object) != o {
// 		t.Fail()
// 	}

// 	type S struct {
// 		Field *js.Object
// 	}
// 	s := S{o}

// 	v := reflect.ValueOf(&s).Elem()
// 	if v.Field(0).Interface().(*js.Object).Get("answer").Int() != 42 {
// 		t.Fail()
// 	}
// 	if v.Field(0).MethodByName("Get").Call([]reflect.Value{reflect.ValueOf("answer")})[0].Interface().(*js.Object).Int() != 42 {
// 		t.Fail()
// 	}
// 	v.Field(0).Set(reflect.ValueOf(js.Global.Call("eval", "({ answer: 100 })")))
// 	if s.Field.Get("answer").Int() != 100 {
// 		t.Fail()
// 	}

// 	if fmt.Sprintf("%+v", s) != "{Field:[object Object]}" {
// 		t.Fail()
// 	}
// }

// func TestNil(t *testing.T) {
// 	type S struct{ X int }
// 	var s *S
// 	if !dummys.Call("isEqual", s, nil).Bool() {
// 		t.Fail()
// 	}

// 	type T struct{ Field *S }
// 	if dummys.Call("testField", T{}) != nil {
// 		t.Fail()
// 	}
// }

// func TestNewArrayBuffer(t *testing.T) {
// 	b := []byte("abcd")
// 	a := js.NewArrayBuffer(b[1:3])
// 	if a.Get("byteLength").Int() != 2 {
// 		t.Fail()
// 	}
// }

// func TestInternalizeExternalizeNull(t *testing.T) {
// 	type S struct {
// 		*js.Object
// 	}
// 	r := js.Global.Call("eval", "(function(f) { return f(null); })").Invoke(func(s S) S {
// 		if s.Object != nil {
// 			t.Fail()
// 		}
// 		return s
// 	})
// 	if r != nil {
// 		t.Fail()
// 	}
// }

// func TestInternalizeExternalizeUndefined(t *testing.T) {
// 	type S struct {
// 		*js.Object
// 	}
// 	r := js.Global.Call("eval", "(function(f) { return f(undefined); })").Invoke(func(s S) S {
// 		if s.Object != js.Undefined {
// 			t.Fail()
// 		}
// 		return s
// 	})
// 	if r != js.Undefined {
// 		t.Fail()
// 	}
// }

// func TestDereference(t *testing.T) {
// 	s := *dummys
// 	p := &s
// 	if p != dummys {
// 		t.Fail()
// 	}
// }

// func TestSurrogatePairs(t *testing.T) {
// 	js.Global.Set("str", "\U0001F600")
// 	str := js.Global.Get("str")
// 	if str.Get("length").Int() != 2 || str.Call("charCodeAt", 0).Int() != 55357 || str.Call("charCodeAt", 1).Int() != 56832 {
// 		t.Fail()
// 	}
// 	if str.String() != "\U0001F600" {
// 		t.Fail()
// 	}
// }

// func TestUint8Array(t *testing.T) {
// 	uint8Array := js.Global.Get("Uint8Array")
// 	if dummys.Call("return", []byte{}).Get("constructor") != uint8Array {
// 		t.Errorf("Empty byte array is not externalized as a Uint8Array")
// 	}
// 	if dummys.Call("return", []byte{0x01}).Get("constructor") != uint8Array {
// 		t.Errorf("Non-empty byte array is not externalized as a Uint8Array")
// 	}
// }

// func TestTypeSwitchJSObject(t *testing.T) {
// 	obj := js.Global.Get("Object").New()
// 	obj.Set("foo", "bar")

// 	want := "bar"

// 	if got := obj.Get("foo").String(); got != want {
// 		t.Errorf("Direct access to *js.Object field gave %q, want %q", got, want)
// 	}

// 	var x interface{} = obj

// 	switch x := x.(type) {
// 	case *js.Object:
// 		if got := x.Get("foo").String(); got != want {
// 			t.Errorf("Value passed through interface and type switch gave %q, want %q", got, want)
// 		}
// 	}

// 	if y, ok := x.(*js.Object); ok {
// 		if got := y.Get("foo").String(); got != want {
// 			t.Errorf("Value passed through interface and type assert gave %q, want %q", got, want)
// 		}
// 	}
// }
=======
>>>>>>> upstream/wasm-wip
