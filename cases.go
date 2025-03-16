// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"reflect"
)

// Cases represents a generic data structure for table-driven testing.
// It is an alias for a slice of cases with the following structure:
//
//	[]struct {
//	  Label string
//	  In1   TIn1
//	   ⋮      ⋮
//	  InN   TInN
//	  Out1  TOut1
//	   ⋮      ⋮
//	  OutM  TInM
//	}
//
// where N and M are arbitrary.
//
// For the common case of testing a numerical function of a single
// argument, this would take the simple form of:
//
//	 []struct {
//	   Label string
//	   In    float64
//	   Out   float64
//	 }{
//		  {"SomeLabel", x, y},
//		  {"AnotherLabel", x1, y1},
//	   ...
//	 }
type Cases interface{}

// parseCases converts cases reflect values and performs basic validation checks.
// If any checks fail, parse panics.
// nc is the number of cases, nf is the number of fields in a case (label + inputs + outputs)
func parseCases(cases Cases) (casesv reflect.Value, ncases, nfields int, err error) {
	casesv = reflect.ValueOf(cases)
	if !casesv.IsValid() {
		err = fmt.Errorf("cases not valid, reflection failed")
		return
	}

	// Ensure cs is a slice of cases.
	kc := casesv.Kind()
	if kc != reflect.Slice {
		err = fmt.Errorf("wrong kind of argument. Got %v, want %v", kc, "slice")
		return
	}

	// Ensure there is at least 1 case.
	ncases = casesv.Len()
	if ncases == 0 {
		err = fmt.Errorf("too few cases. Got 0, want at least 1")
		return
	}

	// Ensure each case is a struct.
	kc = casesv.Type().Elem().Kind()
	if kc != reflect.Struct {
		err = fmt.Errorf("wrong input type. Got []%v, want []%v", kc, "struct")
		return
	}

	// Ensure cases have at least 1 field, for the label.
	nfields = casesv.Index(0).NumField()
	if nfields == 0 {
		err = fmt.Errorf("too few fields in cases. Got 0, want at least 1")
		return
	}

	// Ensure the first field, the label, is a string.
	kfc := casesv.Index(0).Field(0).Kind()
	if kfc != reflect.String {
		err = fmt.Errorf("invalid type for first field. Got %v, want %v", kfc, "string")
		return
	}
	return
}

// name gets the name of a case.
func name(cv reflect.Value) string { return cv.Field(0).String() }

// sliceFrom creates a slice of the fields of the case c.
// The returned slice contains field start to start+n.
//
// For pointer or interface fields, the underlying value
// is used in the output slice.
func sliceFrom(cv reflect.Value, start, n int) []reflect.Value {
	v := make([]reflect.Value, n)
	for i := 0; i < n; i++ {
		v[i] = indirect(cv.Field(start + i))
	}
	return v
}

// indirect returns the value referred to by
// a pointer or interface, or the value itself otherwise.
func indirect(v reflect.Value) reflect.Value {
	if k := v.Kind(); k == reflect.Interface || k == reflect.Ptr || k == reflect.UnsafePointer {
		return v.Elem()
	}
	return v
}
