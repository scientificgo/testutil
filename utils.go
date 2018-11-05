// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import (
	"fmt"
	"reflect"
)

// Func represents an arbitrary function.
type Func interface{}

// Cases represents the data structure holding the cases for running tests
// or benchmarks against. It is an alias for a slice of structs with the
// structure:
//
//  []struct{
//    Label                 string
//    In1, In2, ..., InN    interface{}
//    Out1, Out2, ..., OutM interface{}
//  }
//
// If c is such a slice, then it can be used to run tests as:
//
//  func TestSomeFunc(t *testing.T) { Test(t, digits, c, SomeFunc) }
//
// or benchmarks as:
//
//  func BenchmarkSomeFunc(b *testing.B) { Benchmark(b, c, SomeFunc) }
//
// where SomeFunc is the function being tested or benchmarked with signature
//
//  func SomeFunc(in1 IT1, in2 IT2, ..., inN ITN) (OT1, OT2, ..., OTM)
//
// where IT (OT) stand for input (output) type.
type Cases interface{}

// parseFuncs returns up to two funcv values extracted from fs.
// If only one Func is provided, then f2v is a nil function of
// the same type.
func parseFuncs(fs ...Func) (f1v, f2v funcv) {
	// Get the test function(s) f1v, f2v from f.
	panicIf(len(fs) < 1, "")

	f1v = reflect.ValueOf(fs[0])
	kf := f1v.Kind()
	panicIf(
		kf != reflect.Func,
		"Wrong kind of argument. Got %v, want %v.",
		kf, reflect.Func,
	)

	if len(fs) == 1 {
		f2v = reflect.Zero(f1v.Type())
	} else {
		f2v = reflect.ValueOf(fs[1])
	}
	return
}

// parse converts cases reflect values and performs basic validation checks.
// If any checks fail, parse panics.
func parseCases(cs Cases) (cvs casesv, nc, nf int) {
	cvs = reflect.ValueOf(cs)

	kc := cvs.Kind()
	panicIf(kc != reflect.Slice,
		"Wrong kind of argument. Got %v, want %v.",
		kc, "slice",
	)

	nc = cvs.Len()
	panicIf(nc == 0, "No cases.")

	kc = cvs.Type().Elem().Kind()
	panicIf(kc != reflect.Struct,
		"Wrong input type. Got []%v, want []%v",
		kc, "struct",
	)

	nf = cvs.Index(0).NumField()
	panicIf(nf == 0, "Empty cases.")

	kc = cvs.Index(0).Field(0).Kind()
	panicIf(
		kc != reflect.String,
		"Wrong type for struct label. Got %v, want %v.",
		kc, "string",
	)

	return
}

// funcv is the interface for reflect.ValueOf(f) for Func f.
type funcv interface {
	Call([]reflect.Value) []reflect.Value
	IsNil() bool
	Kind() reflect.Kind
	Type() reflect.Type
}

// casesv is the interface for reflect.ValueOf(c) for Cases c.
type casesv interface {
	Index(int) reflect.Value
	Kind() reflect.Kind
	Len() int
	Type() reflect.Type
}

// casev is the interface for a single case, i.e.,
// casesv.Index(i) for int i.
type casev interface{ Field(int) reflect.Value }

// name gets the name of a case.
func name(cv casev) string { return cv.Field(0).String() }

// slicify creates a slice of the fields of the case c.
// The returned slice contains field start to start+n.
//
// For pointer or interface fields, the underlying value
// is used in the output slice.
func sliceFrom(cv casev, start, n int) []reflect.Value {
	v := make([]reflect.Value, n)
	for i := 0; i < n; i++ {
		v[i] = underlying(cv.Field(start + i))
	}
	return v
}

// underlying returns the underlying value referred to by
// a pointer or interface, or the value itself otherwise.
func underlying(v reflect.Value) reflect.Value {
	if k := v.Kind(); k == reflect.Interface || k == reflect.Ptr {
		return v.Elem()
	}
	return v
}

// panicIf panics if b is true using the format string s and values v.
func panicIf(b bool, s string, v ...interface{}) {
	if b {
		panic(fmt.Sprintf(s, v...))
	}
}
