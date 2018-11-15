// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import (
	"reflect"
	"testing"
)

// Test runs a sub-test for each case in cs using the function(s) in fs.
//
// If one function is given, then the return values for a given case are
// tested against its expected outputs.
//
// Alternatively, if two functions are given then their respective return
// values for a given case are tested against each other; cases do not
// need to contain the expected outputs.
func Test(t *testing.T, digits float64, cs Cases, fs ...Func) {
	cvs, nc, nfc := parseCases(cs)
	f1v, f2v := parseFuncs(fs...)

	nIn := f1v.Type().NumIn()
	nOut := f1v.Type().NumOut()

	validateTestIO(nIn, nOut, nfc, f2v.IsNil())

	for i := 0; i < nc; i++ {
		subtest(t, cvs.Index(i), f1v, f2v, nIn, nOut, digits)
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

// subtest runs a sub-test for a given case.
func subtest(t *testing.T, cv casev, f1v, f2v funcv, nIn, nOut int, digits float64) {
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
			handleSubtest(t, i, ri, oi, digits)
		}
	})
}

// handleSubtest checks whether io and ri are equal and reports if they are not.
func handleSubtest(t *testing.T, i int, oi, ri reflect.Value, digits float64) {
	j, ok := equal(ri, oi, digits)
	if ok {
		return
	}

	if j < 0 {
		t.Errorf("Error: length mismatch between %v-th result and expected output.", i)
	}
	if kind := oi.Kind(); kind == reflect.Slice {
		t.Errorf("Error in results[%v][%v]. Got %v, want %v.",
			i, j, oi.Index(j), ri.Index(j))
	} else if kind == reflect.Struct {
		t.Errorf("Error in results[%v].%v. Got %v, want %v.",
			i, oi.Type().Field(j).Name, oi.Field(j), ri.Field(j))
	} else {
		t.Errorf("Error in results[%v]. Got %v, want %v.", i, oi, ri)
	}
}
