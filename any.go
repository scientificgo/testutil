// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"reflect"
)

// Any returns true if the function f evaluates to true for any argument in xs.
func Any(f Func, xs ...interface{}) (ok bool, err error) {
	// get function and validate it returns bool
	// and that the xs are the right size and type
	fv := reflect.ValueOf(f)
	nArg := len(xs)

	k := fv.Kind()
	if k != reflect.Func {
		err = fmt.Errorf("wrong input type. Got %v, want %v", k, "func")
		return
	}

	k = fv.Type().Out(0).Kind()
	if k != reflect.Bool {
		err = fmt.Errorf("wrong output type. Got %v, want %v", k, reflect.Bool)
		return
	}

	nIn := fv.Type().NumIn()
	if nArg != nIn {
		err = fmt.Errorf("wrong number of input slices. Got %v, want %v", nArg, nIn)
		return
	}

	nOut := fv.Type().NumOut()
	if nOut != 1 {
		err = fmt.Errorf("wrong number of output slices. Got %v, want %v", nOut, 1)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	args := make([]reflect.Value, nArg)
	l := reflect.ValueOf(xs[0]).Len()

	// iterate over the input length
	for i := 0; i < l; i++ {
		// iterate across all inputs and construct the slice for calling f
		for j, x := range xs {
			xv := reflect.ValueOf(x)
			args[j] = indirect(xv.Index(i))
		}
		if fv.Call(args)[0].Interface().(bool) {
			ok = true
			return
		}
	}
	return
}

// All returns true if the function f evaluates to true for all arguments in xs.
func All(f Func, xs ...interface{}) (ok bool, err error) {
	// use all(f) = !any(!f)
	notf := func(y interface{}) bool {
		ys := []reflect.Value{reflect.ValueOf(y)}
		val := reflect.ValueOf(f).Call(ys)[0]
		return !(val.Interface().(bool))
	}
	notOk, err := Any(notf, xs...)
	if err != nil {
		return
	}
	ok = !notOk
	return
}
