// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import "reflect"

// Func represents an arbitrary function with the generic
// signature with N inputs and M outputs of any types, i.e.:
//  func Func(in1 In1, ..., inN InN) (out1 Out1, ..., outM OutM)
type Func interface{}

// funcv is the interface for reflect.ValueOf(f) for Func f.
type funcv interface {
	Call([]reflect.Value) []reflect.Value
	IsNil() bool
	Kind() reflect.Kind
	Type() reflect.Type
}

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
