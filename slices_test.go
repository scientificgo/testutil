// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils_test

import (
	"math"
	"math/cmplx"
	. "scientificgo.org/testutils"
	"testing"
)

const tol = 10

var (
	NaN  = math.NaN()
	Inf  = math.Inf(1)
	cNaN = complex(NaN, NaN)
	cInf = complex(Inf, Inf)
)

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

func TestAllFloat64s(t *testing.T) {
	var f func(*testing.T, float64, []string, func([]float64, func(float64) bool) bool, [][]float64, [](func(float64) bool), []bool)
	GenerateTest(&f)

	labels := []string{"Number", "NotInf1", "Zero"}
	in1 := [][]float64{
		{1, 2, NaN, NaN, NaN},
		{1, 2, 3, 4, Inf},
		{0, 0, 0, 0, 0},
	}

	in2 := []func(float64) bool{
		func(x float64) bool { return math.IsNaN(x) },
		func(x float64) bool { return !math.IsInf(x, 0) },
		func(x float64) bool { return x == 0 },
	}

	out := []bool{false, false, true}

	f(t, tol, labels, AllFloat64s, in1, in2, out)
}

func TestAllComplex128s(t *testing.T) {
	var f func(*testing.T, float64, []string, func([]complex128, func(complex128) bool) bool, [][]complex128, [](func(complex128) bool), []bool)
	GenerateTest(&f)

	labels := []string{"Number", "NotInf1", "Zero"}
	in1 := [][]complex128{
		{1, 2, cNaN, cNaN, cNaN},
		{1, -2i, 3, 4 + 1i, cInf},
		{0, 0, -0, 0, 0},
	}

	in2 := []func(complex128) bool{
		func(x complex128) bool { return cmplx.IsNaN(x) },
		func(x complex128) bool { return !cmplx.IsInf(x) },
		func(x complex128) bool { return x == 0 },
	}

	out := []bool{false, false, true}

	f(t, tol, labels, AllComplex128s, in1, in2, out)
}
