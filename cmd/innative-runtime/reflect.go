// +build ignore

package main

import (
	"reflect"
	"unsafe"
)

// SliceAt returns a view of the memory at p as a slice of elem.
// The elem parameter is the element type of the slice, not the complete slice type.
func SliceAt(elem reflect.Type, p unsafe.Pointer, n int) reflect.Value {
	if p == nil && n == 0 {
		return reflect.Zero(reflect.SliceOf(elem))
	}
	return reflect.NewAt(bigArrayOf(elem), p).Elem().Slice3(0, n, n)
}

// bigArrayOf returns the type of a maximally-sized array of t.
//
// This works around the memory-stranding issue described in
// https://golang.org/issue/13656 by producing only one array type per element
// type (instead of one array type per length).
func bigArrayOf(t reflect.Type) reflect.Type {
	n := ^uintptr(0) / uintptr(t.Size())
	const maxInt = uintptr(^uint(0) >> 1)
	if n > maxInt {
		n = maxInt
	}
	return reflect.ArrayOf(int(n), t)
}
