// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import "reflect"

// Any returns true if the function f evaluates to true for any argument in xs.
func Any(f Func, xs ...interface{}) bool {
	// Get function and validate it returns bool.
	fv := reflect.ValueOf(f)
	nArg := len(xs)
	//validateAny(fv, n)

	k := fv.Kind()
	panicIf(k != reflect.Func, "Wrong input type. Got %v, want %v.", k, "func")

	k = fv.Type().Out(0).Kind()
	panicIf(k != reflect.Bool, "Wrong output type. Got %v, want %v", k, reflect.Bool)

	nIn := fv.Type().NumIn()
	panicIf(nArg != nIn, "Wrong number of input slices. Got %v, want %v.", nArg, nIn)

	nOut := fv.Type().NumOut()
	panicIf(nOut != 1, "Wrong number of output slices. Got %v, want %v.", nOut, 1)

	args := make([]reflect.Value, nArg)
	l := reflect.ValueOf(xs[0]).Len()

	// Iterate over the input length.
	for i := 0; i < l; i++ {
		// Iterate across all inputs and construct the slice for calling f.
		for j, x := range xs {
			xv := reflect.ValueOf(x)
			args[j] = underlying(xv.Index(i))
		}
		if fv.Call(args)[0].Interface().(bool) {
			return true
		}
	}
	return false
}

// All returns true if the function f evaluates to true for all arguments in x.
func All(f Func, x ...interface{}) bool {
	g := func(y interface{}) bool {
		ys := []reflect.Value{reflect.ValueOf(y)}
		return !reflect.ValueOf(f).Call(ys)[0].Interface().(bool)
	}
	return !Any(g, x...)
}
