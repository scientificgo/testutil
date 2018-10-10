// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import (
	"fmt"
	"reflect"
	"testing"
)

// Test is a generic testing tool that can be used to test arbitrary functions
// with any signature to the specified accuracy (in significant digits) when
// provided with suitable cases.
//
// The parameter f can either be a function or an array of exactly two functions.
// If two functions are provided, then their values are compared for the given
// inputs.
//
// The parameter cases must be a slice of structs; the struct must include a string
// as its first field, and all of its variables must be exported. In other words,
// the struct must have the form
//  type myStruct struct {
//    Label string
//    ...
//  }
// otherwise the test will panic.
func Test(t *testing.T, accuracy float64, f interface{}, cases interface{}) {
	// Get the cases and validate their form.
	casesv := reflect.ValueOf(cases)
	if k := casesv.Kind(); k != reflect.Slice {
		panic(fmt.Sprintf("Wrong kind of argument. Got %v, want %v.", k, "slice"))
	}
	if k := casesv.Type().Elem().Kind(); k != reflect.Struct {
		panic(fmt.Sprintf("Wrong input type. Got []%v, want []%v", k, "struct"))
	}
	if k := casesv.Index(0).Field(0).Kind(); k != reflect.String {
		panic(fmt.Sprintf("Wrong struct type for field(0). Got %v, want %v.", k, "string"))
	}

	// Get the test function(s) from f.
	fv := reflect.ValueOf(f)
	var f1, f2 reflect.Value
	if k := fv.Kind(); k == reflect.Func {
		f1 = fv
		f2 = reflect.Zero(f1.Type())
	} else if k == reflect.Array && fv.Type().Elem().Kind() == reflect.Func && fv.Type().Len() == 2 {
		f1 = fv.Index(0)
		f2 = fv.Index(1)
	} else {
		panic(fmt.Sprintf("Wrong input type. Got %v, want %v or [2]%v.", k, "func", "func"))
	}

	nIn := f1.Type().NumIn()
	nOut := f1.Type().NumOut()
	nCases := casesv.Len()
	nFields := casesv.Index(0).NumField()

	// Check the number of IO arguments provided.
	if nFields-1 != nIn+nOut && f2.IsNil() {
		panic(fmt.Sprintf("Wrong number of IO slices. Got %v, want %v.", nFields-1, nIn+nOut))
	} else if nFields-1 < nIn && !f2.IsNil() {
		panic(fmt.Sprintf("Wrong number of IO slices. Got %v, want %v.", nFields-1, nIn))
	}

	// Iterate over each case and run a sub-test.
	for i := 0; i < nCases; i++ {
		c := casesv.Index(i)
		t.Run(c.Field(0).String(), func(t *testing.T) {
			// Iterate across the case to build the inputs
			// (and expected outputs, if required) for the
			// test function(s).
			inputs := make([]reflect.Value, nIn)
			outputs := make([]reflect.Value, nOut)
			for j := 0; j < nIn; j++ {
				inputs[j] = c.Field(1 + j)
			}
			if f2.IsNil() {
				for j := 0; j < nOut; j++ {
					outputs[j] = c.Field(1 + j + nIn)
				}
			} else {
				outputs = f2.Call(inputs)
			}
			results := f1.Call(inputs)
			// Compare results with expected outputs and report
			// any failures with their locations within the data.
			if i, j, ok := approxEqual(results, outputs, accuracy); !ok {
				if k := outputs[i].Kind(); k == reflect.Slice {
					t.Errorf("Error in results[%v][%v]. Got %v, want %v.", i, j, outputs[i].Index(j), results[i].Index(j))
				} else if k == reflect.Struct {
					t.Errorf("Error in results[%v].%v. Got %v, want %v.", i, outputs[i].Type().Field(j).Name, outputs[i].Field(j), results[i].Field(j))
				} else {
					t.Errorf("Error in results[%v]. Got %v, want %v.", i, outputs[i], results[i])
				}
			}
		})
	}
}

// approxEqual recursis through x and y to check they are accurate to the
// desired level and returns the location in x or y where an error occurs,
// as well as the sub-location (index or field name) where appropriate.
func approxEqual(x, y []reflect.Value, acc float64) (int, int, bool) {
	for i := 0; i < len(x); i++ {
		xi, yi := x[i], y[i]
		if k := xi.Kind(); k == reflect.Slice {
			for j := 0; j < xi.Len(); j++ {
				if _, _, ok := approxEqual([]reflect.Value{xi.Index(j)}, []reflect.Value{yi.Index(j)}, acc); !ok {
					return i, j, ok
				}
			}
		} else if k == reflect.Struct {
			for j := 0; j < xi.NumField(); j++ {
				if _, _, ok := approxEqual([]reflect.Value{xi.Field(j)}, []reflect.Value{yi.Field(j)}, acc); !ok {
					return i, j, ok
				}
			}
		} else {
			if !approxEqualVal(xi, yi, acc) {
				return i, 0, false
			}
		}
	}
	return 0, 0, true
}

// approxEqualVal checks whether the values x and y are equal to
// within the specified accuracy in significant digits.
func approxEqualVal(x, y reflect.Value, acc float64) bool {
	switch x.Interface().(type) {
	case float64:
		if !EqualFloat64(x.Interface().(float64), y.Interface().(float64), acc) {
			return false
		}
	case complex128:
		if !EqualComplex128(x.Interface().(complex128), y.Interface().(complex128), acc) {
			return false
		}
	default: // int, bool, string, byte
		if !reflect.DeepEqual(x.Interface(), y.Interface()) {
			return false
		}
	}
	return true
}
