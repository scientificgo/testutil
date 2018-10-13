// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import (
	"reflect"
)

// Any returns true if the function f evaluates to true
// for any argument in x.
func Any(f interface{}, x ...interface{}) bool {
	// Get function and validate it returns bool.
	fv := reflect.ValueOf(f)
	n := len(x)
	validateAny(fv, n)

	args := make([]reflect.Value, n)
	l := reflect.ValueOf(x[0]).Len()
	for i := 0; i < l; i++ {
		for j, xj := range x {
			xjv := reflect.ValueOf(xj)
			args[j] = underlyingValue(xjv.Index(i))
		}

		if fv.Call(args)[0].Interface().(bool) {
			return true
		}
	}
	return false
}

// All returns true if the function f evaluates to true
// for all arguments in x.
func All(f interface{}, x ...interface{}) bool {
	notf := func(y interface{}) bool {
		return !reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(y)})[0].Interface().(bool)
	}
	return !Any(notf, x...)
}
