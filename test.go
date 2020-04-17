// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import (
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
//
// For example, to test a function called SquareRoot on some cases, and
// against the standard library function math.Sqrt, where SquareRoot,
// cases and the tolerance tol are defined as
//
//  func SquareRoot(x float64) float64 { ... }
//
//  var cases = []struct{
//    Label     string
//    In1, Out1 interface{}
//  }{
//    {"Case1", 1., 1.},
//    {"Case2", 4., 2.},
//    {"Case3", -1., math.NaN()},
//  }
//
//  var tol = 1.e-10
//
// then the test functions would be
//
//  func TestSquareRoot(t *testing.T) {
//    testutil.Test(t, tol, cases, SquareRoot)
//  }
//
//  func TestSquareRootvsSqrt(t *testing.T) {
//    testutil.Test(t, tol, cases, SquareRoot, math.Sqrt)
//  }
func Test(t *testing.T, tolerance interface{}, cs Cases, fs ...Func) {
	tol := tolerance.(float64)
	cvs, nc, nfc := parseCases(cs)
	f1v, f2v := parseFuncs(fs...)

	nIn := f1v.Type().NumIn()
	nOut := f1v.Type().NumOut()

	validateTestIO(nIn, nOut, nfc, f2v.IsNil())

	for i := 0; i < nc; i++ {
		subtest(t, cvs.Index(i), f1v, f2v, nIn, nOut, tol)
	}
}

// validateTestIO panics if the provided arguments are inconsistent.
func validateTestIO(nIn, nOut, nfc int, f2vIsNil bool) {
	panicIf(
		nfc-1 != nIn+nOut && f2vIsNil,
		"Wrong number of input/output slices. Got %v, want %v.",
		nfc-1, nIn+nOut,
	)
	panicIf(
		nfc-1 != nIn && !f2vIsNil,
		"Wrong number of input slices. Got %v, want %v.",
		nfc-1, nIn,
	)
}

// subtest runs a subtest for a case.
func subtest(t *testing.T, cv casev, f1v, f2v funcv, nIn, nOut int, tol float64) {
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
			handleSubtest(t, i, ri, oi, tol)
		}
	})
}

// handleSubtest handles the output of the comparison between ri and oi.
func handleSubtest(t *testing.T, i int, ri, oi reflect.Value, tol float64) {
	if res := equal(ri, oi, tol); !res.Ok {
		j := res.Position
		if j < 0 {
			t.Errorf("Error: length mismatch between %v-th result and expected output.", i)
		}
		switch kind := oi.Kind(); {
		case kind == reflect.Slice:
			if !res.Numerical {
				t.Errorf("Error in results[%v][%v]. Got %v, want %v.",
					i, j, ri.Index(j), oi.Index(j))
			} else {
				t.Errorf("Error in results[%v][%v]. Got %v, want %v. (%v)",
					i, j, ri.Index(j), oi.Index(j), res.RelativeError)
			}

		case kind == reflect.Struct:
			if !res.Numerical {
				t.Errorf("Error in results[%v].%v. Got %v, want %v.",
					i, oi.Type().Field(j).Name, ri.Field(j), oi.Field(j))
			} else {
				t.Errorf("Error in results[%v].%v. Got %v, want %v. (%v)",
					i, oi.Type().Field(j).Name, ri.Field(j), oi.Field(j), res.RelativeError)
			}

		default:
			if !res.Numerical {
				t.Errorf("Error in results[%v]. Got %v, want %v.", i, ri, oi)
			} else {
				t.Errorf("Error in results[%v]. Got %v, want %v. (%v)", i, ri, oi, res.RelativeError)
			}
		}
	}
}
