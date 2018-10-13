// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import (
	"fmt"
	"reflect"
)

func panicf(b bool, s string, v ...interface{}) {
	if b {
		panic(fmt.Sprintf(s, v...))
	}
}

func validateAny(fv reflect.Value, nArg int) {
	var k reflect.Kind

	k = fv.Kind()
	panicf(k != reflect.Func, "Wrong input type. Got %v, want %v.", k, "func")

	nIn := fv.Type().NumIn()
	panicf(nArg != nIn, "Wrong number of input slices. Got %v, want %v.", nArg, nIn)

	nOut := fv.Type().NumOut()
	panicf(nOut != 1, "Wrong number of output slices. Got %v, want %v.", nOut, 1)

	k = fv.Type().Out(0).Kind()
	panicf(k != reflect.Bool, "Wrong output type. Got %v, want %v", k, reflect.Bool)
}

func validateTest(f, cases interface{}) (casesv, f1, f2 reflect.Value, nCases, nIn, nOut int) {
	casesv, nCases, nField := validateCases(cases)

	// Get the test function(s) f1, f2 from f.
	fv := reflect.ValueOf(f)
	if k := fv.Kind(); k == reflect.Func {
		f1 = fv
		f2 = reflect.Zero(f1.Type())
	} else if k == reflect.Array && fv.Type().Elem().Kind() == reflect.Func && fv.Type().Len() == 2 {
		f1 = fv.Index(0)
		f2 = fv.Index(1)
	} else {
		panicf(true, "Wrong input type. Got %v, want %v or [2]%v.", k, reflect.Func, reflect.Func)
	}

	nIn = f1.Type().NumIn()
	nOut = f1.Type().NumOut()

	panicf(nField-1 != nIn+nOut && f2.IsNil(), "Wrong number of input/output slices. Got %v, want %v.", nField-1, nIn+nOut)
	panicf(nField-1 != nIn && !f2.IsNil(), "Wrong number of input slices. Got %v, want %v.", nField-1, nIn)
	return
}

func validateBenchmark(f, cases interface{}) (casesv, f1 reflect.Value, nCases, nIn int) {
	casesv, nCases, nField := validateCases(cases)

	f1 = reflect.ValueOf(f)
	k := f1.Kind()
	panicf(k != reflect.Func, "Wrong input type. Got %v, want %v.", k, reflect.Func)

	nIn = f1.Type().NumIn()
	panicf(nField-1 != nIn, "Wrong number of input slices. Got %v, want %v.", nField-1, nIn)
	return
}

func validateCases(cases interface{}) (casesv reflect.Value, nCases, nField int) {
	casesv = reflect.ValueOf(cases)
	var k reflect.Kind

	k = casesv.Kind()
	panicf(k != reflect.Slice, "Wrong kind of argument. Got %v, want %v.", k, "slice")

	nCases = casesv.Len()
	panicf(nCases == 0, "No cases.")

	k = casesv.Type().Elem().Kind()
	panicf(k != reflect.Struct, "Wrong input type. Got []%v, want []%v", k, "struct")

	nField = casesv.Index(0).NumField()
	panicf(nField == 0, "Empty cases.")

	k = casesv.Index(0).Field(0).Kind()
	panicf(k != reflect.String, "Wrong type for struct label. Got %v, want %v.", k, "string")
	return
}
