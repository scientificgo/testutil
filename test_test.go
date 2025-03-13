// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil_test

import (
	"math"
	"testing"

	. "github.com/scientificgo/testutil"
)

func TestTest_Slices(t *testing.T) {
	tol := 0.001
	cases := []struct {
		Label          string
		In, Out1, Out2 interface{}
	}{
		{"", []float64{1, 2}, []float64{1, 2}, []string{"a", "b"}},
		{"", []float64{1, 2.00001}, []float64{1, 2}, []string{"a", "b"}},
	}
	f := func(in []float64) ([]float64, []string) { return in, []string{"a", "b"} }
	Test(t, tol, cases, f)
}

func TestTest_Structs(t *testing.T) {
	tol := 1
	cases := []struct {
		Label    string
		In, Out1 interface{}
	}{
		{"", struct{ A, B int }{1, 2}, struct{ A, B int }{1, 2}},
		{"", struct{ A, B int }{1, 3}, struct{ A, B int }{1, 2}},

		{"", struct{ A, B bool }{true, true}, struct{ A, B bool }{true, true}},
	}
	f := func(in struct{ A, B int }) struct{ A, B int } { return in }
	g := func(in struct{ A, B bool }) struct{ A, B bool } { return in }

	Test(t, tol, cases[:2], f)
	Test(t, tol, cases[2:], g)
}

func TestTest_Maps(t *testing.T) {
	tol := 1
	cases := []struct {
		Label    string
		In, Out1 interface{}
	}{
		{"", map[string]int{"a": 1, "b": 3}, map[string]int{"a": 1, "b": 2}},
		{"", map[string]int{"a": 1, "b": 3}, map[string]int{"a": 1, "b": 2}},

		{"", map[string]string{"a": "1", "b": "2"}, map[string]string{"a": "1", "b": "2"}},

		{"",
			map[int](map[string]string){1: map[string]string{"a": "1", "b": "2"}},
			map[int](map[string]string){1: map[string]string{"a": "1", "b": "2"}},
		},
	}
	f := func(in map[string]int) map[string]int { return in }
	g := func(in map[string]string) map[string]string { return in }
	h := func(in map[int](map[string]string)) map[int](map[string]string) { return in }

	Test(t, tol, cases[:2], f)
	Test(t, nil, cases[2:3], g)
	Test(t, nil, cases[3:], h)
}

func TestTest_Funcs(t *testing.T) {
	tol := 0.1

	cases := []struct {
		Label         string
		In1, In2, Out float64
	}{
		{"", 0, 0, 0},
		{"", 1, 1, math.Sqrt2},
	}

	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}

	Test(t, tol, cases, math.Hypot)
	Test(t, tol, cases, math.Hypot, hypot)
}

func TestTest_Default(t *testing.T) {
	chani64 := make(chan int64)
	chanf64 := make(chan float64)

	cases := []struct {
		Label         string
		In1, In2, Out interface{}
	}{
		{"", chani64, chani64, true},
		{"", chanf64, chani64, false},
		{"", 1 == 0, 1 == 1*1, false},
		{"", "Ben", "Jerry", false},
		{"", complex128(1.0), complex128(2.0), complex128(300.0)},
	}

	f := func(x, y interface{}) bool {
		res := Equal(x, y, nil)
		return res.Ok
	}
	g := func(x, y complex128) complex128 { return x + y }

	Test(t, nil, cases[:4], f)
	Test(t, -inf, cases[4:], g)
}

// The following tests fail by design and are used to
// check the error messages produced are as expected.

// func TestTest_WrongLengthError(t *testing.T) {
// 	cases := []struct {
// 		Label          string
// 		In, Out1, Out2 []float64
// 	}{{"", []float64{1, 2}, []float64{1, 2, 3}, []float64{1, 2}}, // length mismatch in out[0]
// 		{"", []float64{1, 2}, []float64{1, 2}, []float64{1, 2, 3}}, // length mismatch in out[1]
// 		{"", []float64{1, 2, 3}, []float64{1, 2}, []float64{1, 2}}, // length mismatch in out[0] and out[1]
// 	}
// 	f := func(in []float64) ([]float64, []float64) { return in, in }
// 	Test(t, nil, cases, f)
// }

// func TestTest_SliceErrors(t *testing.T) {
// 	cases := []struct {
// 		Label          string
// 		In, Out1, Out2 interface{}
// 	}{
// 		{"", []float64{1, 2}, []float64{1, 2}, []string{"a", "c"}},   // value of out[1]
// 		{"", []float64{1, 2.1}, []float64{1, 2}, []string{"a", "b"}}, // value of out[0] (numerical)

// 	}
// 	f := func(in []float64) ([]float64, []string) { return in, []string{"a", "b"} }
// 	Test(t, nil, cases, f)
// }

// func TestTest_StructErrors(t *testing.T) {
// 	cases := []struct {
// 		Label    string
// 		In, Out1 interface{}
// 	}{
// 		{"", struct{ A, B int }{1, 2}, struct{ A, C int }{1, 2}}, // missing field C
// 		{"", struct{ A, B int }{1, 3}, struct{ A, B int }{1, 2}}, // value of B (numerical)

// 		{"", struct{ A, B bool }{true, true}, struct{ A, B bool }{true, false}}, // value of B
// 	}
// 	f := func(in struct{ A, B int }) struct{ A, B int } { return in }
// 	g := func(in struct{ A, B bool }) struct{ A, B bool } { return in }

// 	Test(t, nil, cases[:2], f)
// 	Test(t, nil, cases[2:], g)
// }

// func TestTest_MapErrors(t *testing.T) {
// 	cases := []struct {
// 		Label    string
// 		In, Out1 interface{}
// 	}{
// 		{"", map[string]int{"a": 1, "c": 2}, map[string]int{"a": 1, "b": 2}}, // missing key "b"
// 		{"", map[string]int{"a": 1, "b": 3}, map[string]int{"a": 1, "b": 2}}, // value of "b" (numerical)

// 		{"", map[string]string{"a": "1", "b": "two"}, map[string]string{"1": "1", "2": "2"}},

// 		{"",
// 			map[int](map[string]string){1: map[string]string{"a": "1", "b": "two"}},
// 			map[int](map[string]string){1: map[string]string{"a": "1", "b": "2"}},
// 		}, // value of "b" for nested map
// 	}
// 	f := func(in map[string]int) map[string]int { return in }
// 	g := func(in map[string]string) map[string]string { return in }
// 	h := func(in map[int](map[string]string)) map[int](map[string]string) { return in }

// 	Test(t, nil, cases[:2], f)
// 	Test(t, nil, cases[2:3], g)
// 	Test(t, nil, cases[3:], h)
// }

// func TestTest_DefaultErrors(t *testing.T) {
// 	cases := []struct {
// 		Label         string
// 		In1, In2, Out interface{}
// 	}{
// 		{"", testing.T{}, testing.B{}, true}, // should be false
// 		{"", complex128(1.0), complex128(2.0), complex128(2.5)},
// 	}

// 	f := func(x, y interface{}) bool {
// 		res := Equal(x, y, nil)
// 		return res.Ok
// 	}
// 	g := func(x, y complex128) complex128 { return x + y }

// 	Test(t, nil, cases[:1], f)
// 	Test(t, nil, cases[1:], g)
// }
