// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import "fmt"
import "testing"
import "reflect"

// GenerateTest initializes a test function, compatible with the standard
// package testing, given a pointer fptr. The pointer fptr should point to
// a function with a signature following one of the below two templates:
//
//  func(*testing.T, float64, []string, func(I1, I2, ..., In) (O1, O2, ..., Om), []I1, []I2, ..., []In, []O1, []O2, ..., []Om)
//
//  func(*testing.T, float64, []string, [2]func(I1, I2, ..., In) (O1, O2, ..., Om), []I1, []I2, ..., []In)
//
// where I1, I2, etc., indicate input types and O1, O2, etc., indicate output
// types. The first three arguments are compulsory; they are for the special
// testing.T variable, the numerical tolerance (i.e. required accuracy in
// significant digits) and sub-test labels.
//
// If two functions are passed then the output slices can be omitted
// (i.e. the second template); in this case, the outputs of the functions
// are compared element-wise.
//
//
// This makes it simple to define custom tests with minimal repitition.
// For example, to test the function
//  func MyFunc(float64, []int) []float64
// against known outputs, the complete test function would be:
//  func TestMyFunc(t *testing.T) {
//    var testf func(*testing.T, float64, []string, func(float64, []int) []float64, []float64, [][]int, [][]float64)
//    GenerateTest(&testf)
//    // Define inputs for (sub-)tests.
//    const tol = 8 // required accuracy in significant digits
//    labels := []string{...}
//    inputs1 := []float64{...}
//    inputs2 := [][]int{...}
//    outputs := [][]float64{...}
//    // Run the (sub-)tests.
//    testf(t, tol, labels, MyFunc, inputs1, inputs2, outputs)
//  }
func GenerateTest(fptr interface{}) { generate(fptr, test) }

func generate(fptr interface{}, template func([]reflect.Value) []reflect.Value) {
	f := reflect.ValueOf(fptr).Elem()
	v := reflect.MakeFunc(f.Type(), template)
	f.Set(v)
}

// test is a generic function for numerical testing, with:
// in[0]: *testing.T
// in[1]: float64 (tolerance in significant digits)
// in[2]: []string (labels)
// in[3]: func or []func (function(s) to test)
// in[4:]: slices (1 per input parameter and output value)
func test(in []reflect.Value) []reflect.Value {
	t := in[0].Interface().(*testing.T)
	tol := in[1].Interface().(float64)
	labels := in[2].Interface().([]string)

	fpos := 3
	var f, g reflect.Value
	if k := in[fpos].Kind(); k == reflect.Slice {
		if in[fpos].Type().Elem().Kind() != reflect.Func {

		}
		f = in[fpos].Index(0)
		g = in[fpos].Index(1)
	} else if k == reflect.Func {
		f = in[fpos]
		g = reflect.Zero(f.Type())
	} else {
		panic(fmt.Sprintf("Wrong input type. Got %v, want %v.", k.String(), "Func"))
	}
	fpos++

	nIn := f.Type().NumIn()
	nOut := f.Type().NumOut()
	nCases := in[fpos].Len()

	for i := fpos; i < len(in); i++ {
		if k := in[i].Kind(); k != reflect.Slice {
			panic(fmt.Sprintf("Wrong input/output type. Got %v, want %v.", k.String(), "Slice"))
		}
		if in[i].Len() != nCases {
			panic(fmt.Sprintf("Input/output lengths do not match. Got %v, want %v.", in[i].Len(), nCases))
		}
	}

	if len(in)-fpos != nIn+nOut && g.IsNil() {
		panic(fmt.Sprintf("Wrong number of input/output slices. Got %v, want %v.", len(in)-fpos, nIn+nOut))
	}

	inputSlices := in[fpos : fpos+nIn]
	outputSlices := in[fpos+nIn:]

	for i := 0; i < nCases; i++ {
		t.Run(labels[i], func(t *testing.T) {
			inputs := make([]reflect.Value, nIn)
			outputs := make([]reflect.Value, nOut)
			for j := 0; j < nIn; j++ {
				inputs[j] = inputSlices[j].Index(i)
			}
			if g.IsNil() {
				for j := 0; j < nOut; j++ {
					outputs[j] = outputSlices[j].Index(i)
				}
			} else {
				outputs = g.Call(inputs)
			}

			results := f.Call(inputs)
			if where, pos, ok := approxEqual(results, outputs, tol); !ok {
				if k := outputs[where].Kind(); k == reflect.Slice {
					t.Errorf("Error (%.0vsd) in results[%v][%v]. Got %v, wanted %v.", tol, where, pos, outputs[where].Index(pos), results[where].Index(pos))
				} else if k == reflect.Struct {
					t.Errorf("Error (%.0vsd) in results[%v].%v. Got %v, wanted %v.", tol, where, outputs[where].Type().Field(pos).Name, outputs[where].Field(pos), results[where].Field(pos))
				} else {
					t.Errorf("Error (%.0vsd) in results[%v]. Got %v, wanted %v.", tol, where, outputs[where], results[where])
				}
			}
		})
	}

	return nil
}

func approxEqual(x, y []reflect.Value, tol float64) (int, int, bool) {
	for i := 0; i < len(x); i++ {
		xi, yi := x[i], y[i]
		if k := xi.Kind(); k == reflect.Slice {
			for j := 0; j < xi.Len(); j++ {
				if _, _, ok := approxEqual([]reflect.Value{xi.Index(j)}, []reflect.Value{yi.Index(j)}, tol); !ok {
					return i, j, ok
				}
			}
		} else if k == reflect.Struct {
			for j := 0; j < xi.NumField(); j++ {
				if _, _, ok := approxEqual([]reflect.Value{xi.Field(j)}, []reflect.Value{yi.Field(j)}, tol); !ok {
					return i, j, ok
				}
			}
		} else {
			if !approxEqualVal(xi, yi, tol) {
				return i, 0, false
			}
		}
	}
	return 0, 0, true
}

func approxEqualVal(x, y reflect.Value, tol float64) bool {
	switch x.Interface().(type) {
	case float64:
		if !EqualFloat64(x.Interface().(float64), y.Interface().(float64), tol) {
			return false
		}
	case complex128:
		if !EqualComplex128(x.Interface().(complex128), y.Interface().(complex128), tol) {
			return false
		}
	default:
		if !reflect.DeepEqual(x.Interface(), y.Interface()) {
			return false
		}
	}
	return true
}
