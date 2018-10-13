// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import (
	"reflect"
	"testing"
)

// Test is a generic testing tool that tests the outputs of an
// arbitrary function f to the specified number of significant digits
// (for float or complex types).
//
// If f is an array [2]func(...)(...) then the respective outputs
// are tested against each other; otherwise the expected output
// should be specified.
//
// cases must be a slice of structs that:
// (i) have a first field that is a string;
// (ii) export all fields.
func Test(t *testing.T, digits float64, f, cases interface{}) {
	casesv, f1, f2, nCases, nIn, nOut := validateTest(f, cases)

	// Iterate over each case and run a sub-test.
	for i := 0; i < nCases; i++ {
		c := casesv.Index(i)
		t.Run(c.Field(0).String(), func(t *testing.T) {
			inputs := sliced(c, 1, nIn)
			var outputs []reflect.Value
			if f2.IsNil() {
				outputs = sliced(c, 1+nIn, nOut)
			} else {
				outputs = f2.Call(inputs)
			}
			results := f1.Call(inputs)
			for j := 0; j < nOut; j++ {
				execSubtest(t, results[j], outputs[j], j, digits)
			}
		})
	}
}

func execSubtest(t *testing.T, result, output reflect.Value, i int, digits float64) {
	if j, ok := equal(result, output, digits); !ok {
		if i < 0 {
			t.Errorf("Error: length mismatch between results and expected outputs.")
		}
		if j < 0 {
			t.Errorf("Error: length mismatch between %v-th result and expected output.", i)
		}

		if k := output.Kind(); k == reflect.Slice {
			t.Errorf("Error in results[%v][%v]. Got %v, want %v.", i, j, output.Index(j), result.Index(j))
		} else if k == reflect.Struct {
			t.Errorf("Error in results[%v].%v. Got %v, want %v.", i, output.Type().Field(j).Name, output.Field(j), result.Field(j))
		} else {
			t.Errorf("Error in results[%v]. Got %v, want %v.", i, output, result)
		}
	}
}
