// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"reflect"
)

// Func represents an arbitrary function with a generic
// signature of N inputs and M outputs of any types, i.e.:
//  func Func(in1 In1, ..., inN InN) (out1 Out1, ..., outM OutM)
type Func interface{}

// parseFuncs parses a slice of Funcs and returns the underlying
// functions as reflect.Values. It returns an error if fewer than 1 or
// more than 2 functions are given, or if a non-function argument is provided.
func parseFuncs(funcs ...Func) (func1v, func2v reflect.Value, err error) {
	l := len(funcs)
	if l < 1 || l > 2 {
		err = fmt.Errorf("wrong number of functions. Got %v, want 1 or 2", l)
		return
	}

	func1v = reflect.ValueOf(funcs[0])
	if !func1v.IsValid() {
		err = fmt.Errorf("invalid functions, reflection failed")
		return
	}

	kf := func1v.Kind()
	if kf != reflect.Func {
		err = fmt.Errorf("wrong kind of argument. Got %v, want %v", kf, reflect.Func)
		return
	}

	if l == 1 {
		func2v = reflect.Zero(func1v.Type()) // f2v cast as a Zero value of the type of f1v
	} else {
		func2v = reflect.ValueOf(funcs[1])
	}
	return
}
