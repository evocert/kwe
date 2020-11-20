package es51

import (
	"runtime"
	"testing"
)

func TestWeakMapExpiry(t *testing.T) {
	vm := New()
	_, err := vm.RunString(`
	var m = new WeakMap();
	var key = {};
	m.set(key, true);
	if (!m.has(key)) {
		throw new Error("has");
	}
	if (m.get(key) !== true) {
		throw new Error("value does not match");
	}
	key = undefined;
	`)
	if err != nil {
		t.Fatal(err)
	}
	runtime.GC()
	runtime.GC()
	vm.RunString("true") // this will trigger dead keys removal
	wmo := vm.Get("m").ToObject(vm).self.(*weakMapObject)
	l := len(wmo.m.data)
	if l > 0 {
		t.Fatal("Object has not been removed")
	}
}
