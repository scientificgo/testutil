// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"reflect"
	"testing"
)

// Test is a generic case-driven testing function that accepts a
// slice of cases, a numerical tolerance and either 1 or 2 functions
// to be tested. A sub-test is run for each case.
//
// If 1 function is provided, then its output is tested against
// the outputs provided in each case.
//
// If 2 functions are provided, then their respective outputs are
// compared, using the inputs provided in each case.
func Test(t *testing.T, tolerance interface{}, cases Cases, funcs ...Func) {
	tol := validateTolerance(tolerance)
	cvs, nc, nfc, err := parseCases(cases)
	if err != nil {
		t.Fatal(err)
	}
	f1v, f2v, err := parseFuncs(funcs...)
	if err != nil {
		t.Fatal(err)
	}
	nIn := f1v.Type().NumIn()
	nOut := f1v.Type().NumOut()

	switch f2v.IsNil() {
	case true: // 1 func
		if nfc-1 != nIn+nOut {
			t.Fatalf("wrong number of input/output slices. Got %v, want %v", nfc-1, nIn+nOut)
		}
	case false: // 2 funcs
		if nfc-1 != nIn+nOut && nfc-1 != nIn { // outputs are optional with 2 funcs
			t.Fatalf("wrong number of input slices. Got %v, want %v", nfc-1, nIn)
		}
	}

	for i := 0; i < nc; i++ {
		subtest(t, cvs.Index(i), f1v, f2v, nIn, nOut, tol)
	}
}

// subtest runs a subtest for a case.
func subtest(t *testing.T, cv, f1v, f2v reflect.Value, nIn, nOut int, tol float64) {
	t.Run(name(cv), func(t *testing.T) {
		var in, out, res []reflect.Value

		in = sliceFrom(cv, 1, nIn)
		if f2v.IsNil() {
			out = sliceFrom(cv, 1+nIn, nOut)
		} else {
			out = f2v.Call(in)
		}
		res = f1v.Call(in)

		for i := 0; i < nOut; i++ {
			ri := res[i]
			oi := out[i]
			if err := handleSubtest(i, ri, oi, tol); err != nil {
				t.Error(err)
			}
		}
	})
}

// handleSubtest returns an error if a subtest fails.
func handleSubtest(i int, ri, oi reflect.Value, tol float64) (err error) {
	res := equal(ri, oi, tol)
	if res.Ok {
		return
	}
	if res.LengthMismatch {
		err = fmt.Errorf("[%v]: Length mismatch", i)
		return
	}
	if res.MissingValue {
		missing := res.Position
		switch kind := oi.Kind(); kind {
		case reflect.Struct:
			err = fmt.Errorf("[%v]: Missing struct field %v", i, oi.Type().Field(missing).Name)
		case reflect.Map:
			err = fmt.Errorf("[%v]: Missing key %v", i, oi.MapKeys()[missing])
		default:
			err = fmt.Errorf("[%v]: Should never reach here", i)
		}
		return
	}

	pos := res.Position

	switch res.Numerical {
	case true:
		switch kind := oi.Kind(); kind {
		case reflect.Struct:
			err = fmt.Errorf("[%v].%v: Got %v, want %v (δ=%v)", i, oi.Type().Field(pos).Name,
				ri.Field(pos), oi.Field(pos), res.RelativeError)
		case reflect.Map:
			key := oi.MapKeys()[pos]
			err = fmt.Errorf("[%v][%v]: Got %v, want %v (δ=%v)", i, key,
				ri.MapIndex(key), oi.MapIndex(key), res.RelativeError)
		case reflect.Array, reflect.Slice:
			err = fmt.Errorf("[%v][%v]: Got %v, want %v (δ=%v)", i, pos,
				ri.Index(pos), oi.Index(pos), res.RelativeError)
		default:
			err = fmt.Errorf("[%v]: Got %v, want %v (δ=%v)", i, ri, oi, res.RelativeError)
		}
		return

	default:
		switch kind := oi.Kind(); kind {
		case reflect.Struct:
			err = fmt.Errorf("[%v].%v: Got %v, want %v", i, oi.Type().Field(pos).Name,
				ri.Field(pos), oi.Field(pos))
		case reflect.Map:
			key := oi.MapKeys()[pos]
			err = fmt.Errorf("[%v][%v]: Got %v, want %v", i, key,
				ri.MapIndex(key), oi.MapIndex(key))
		case reflect.Array, reflect.Slice:
			err = fmt.Errorf("[%v][%v]: Got %v, want %v", i, pos,
				ri.Index(pos), oi.Index(pos))
		default:
			err = fmt.Errorf("[%v]: Got %v, want %v", i, ri, oi)
		}
		return
	}
}
