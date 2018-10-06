// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package utils_test

import (
	. "github.com/scientificgo/utils"
	"math"
	"testing"
)

const tol = 10

var (
	NaN  = math.NaN()
	Inf  = math.Inf(1)
	cNaN = complex(NaN, NaN)
	cInf = complex(Inf, Inf)
)

func TestEqualFloat64(t *testing.T) {
	cases := []struct {
		name      string
		x, y, tol float64
		out       bool
	}{
		{"NaN", NaN, NaN, 10, true},
		{"Zero", 0, 1e-9, 10, false},
		{"Equal", math.Pi, math.Pi, 20, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := EqualFloat64(c.x, c.y, c.tol)
			if res != c.out {
				t.Errorf("EqualFloat64(%v, %v, %v) = %v, want %v", c.x, c.y, c.tol, res, c.out)
			}
		})
	}
}

func TestEqualFloat64s(t *testing.T) {
	cases := []struct {
		x, y []float64
		tol  float64
		out  bool
	}{
		{[]float64{1, 2}, []float64{1, 2, 3}, 10, false},
		{[]float64{0, 0}, []float64{0, 0}, 10, true},
		{[]float64{0, NaN}, []float64{0, NaN}, 10, true},
		{[]float64{0.9999, 1.1121}, []float64{0.999, 1.112}, 5, false},
		{[]float64{0.99, 1.1}, []float64{0.99999, 1.1}, 2, true},
		{[]float64{0.99, 1, 0}, []float64{1, 0.99}, 2, false},
	}
	for _, c := range cases {
		res := EqualFloat64s(c.x, c.y, c.tol)
		if res != c.out {
			t.Errorf("EqualFloat64s(%v, %v, %v) = %v, want %v", c.x, c.y, c.tol, res, c.out)
		}
	}
}

func TestEqualComplex128s(t *testing.T) {
	cases := []struct {
		x, y []complex128
		tol  float64
		out  bool
	}{
		{[]complex128{1 + 1i, 2 + 1i, 3 - 3i}, []complex128{1 + 1i, 2 + 1i}, 10, false},
		{[]complex128{1 + 1i, 2 + 1i, 3 - 3i}, []complex128{1 + 1i, 2 + 1i, 3.01 - 3i}, 10, false},
	}
	for _, c := range cases {
		res := EqualComplex128s(c.x, c.y, c.tol)
		if res != c.out {
			t.Errorf("EqualComplex128s(%v, %v, %v) = %v, want %v", c.x, c.y, c.tol, res, c.out)
		}
	}
}

func TestAnyFloat64s(t *testing.T) {
	cases := []struct {
		name  string
		vals  []float64
		testf func(float64) bool
		out   bool
	}{
		{"HasNaN", []float64{1, 2, 3, 4, NaN}, func(y float64) bool { return math.IsNaN(y) }, true},
		{"HasNaN", []float64{1, 2, 3, 4, 5}, func(y float64) bool { return math.IsNaN(y) }, false},
		{"HasInf", []float64{1, 2, 3, 4, Inf}, func(y float64) bool { return math.IsInf(y, 0) }, true},
		{"HasInf", []float64{1, 2, 3, 4, 5}, func(y float64) bool { return math.IsInf(y, 0) }, false},
		{"HasZero", []float64{1, 2, 3, 4, 0}, func(y float64) bool { return y == 0 }, true},
		{"HasZero", []float64{1, 2, 3, 4, 5}, func(y float64) bool { return y == 0 }, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := AnyFloat64s(c.vals, c.testf)
			if res != c.out {
				t.Errorf("%v(%v) = %v, want %v", c.name, c.vals, res, c.out)
			}
		})
	}
}

func TestAllFloat64s(t *testing.T) {
	cases := []struct {
		name  string
		vals  []float64
		testf func(float64) bool
		out   bool
	}{
		{"IsNaN", []float64{NaN, NaN, NaN}, func(y float64) bool { return math.IsNaN(y) }, true},
		{"HasNaN", []float64{NaN, NaN, NaN, 3}, func(y float64) bool { return math.IsNaN(y) }, false},
		{"HasZero", []float64{0, 0, 0, 0, 0}, func(y float64) bool { return y == 0 }, true},
		{"HasZero", []float64{1, 0, 3, 0, 5}, func(y float64) bool { return y == 0 }, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := AllFloat64s(c.vals, c.testf)
			if res != c.out {
				t.Errorf("%v(%v) = %v, want %v", c.name, c.vals, res, c.out)
			}
		})
	}
}

func TestAnyComplex128s(t *testing.T) {
	cases := []struct {
		name  string
		vals  []complex128
		testf func(complex128) bool
		out   bool
	}{
		{"HasZero", []complex128{1, 2, 3, 4, 0}, func(y complex128) bool { return y == 0 }, true},
		{"HasZero", []complex128{1, 2, 3, 4, 5}, func(y complex128) bool { return y == 0 }, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := AnyComplex128s(c.vals, c.testf)
			if res != c.out {
				t.Errorf("%v(%v) = %v, want %v", c.name, c.vals, res, c.out)
			}
		})
	}
}

func TestAllComplex128s(t *testing.T) {
	cases := []struct {
		name  string
		vals  []complex128
		testf func(complex128) bool
		out   bool
	}{
		{"HasZero", []complex128{0, 0, 0, 0, 0}, func(y complex128) bool { return y == 0 }, true},
		{"HasZero", []complex128{1, 0, 3, 0, 5}, func(y complex128) bool { return y == 0 }, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := AllComplex128s(c.vals, c.testf)
			if res != c.out {
				t.Errorf("%v(%v) = %v, want %v", c.name, c.vals, res, c.out)
			}
		})
	}
}
