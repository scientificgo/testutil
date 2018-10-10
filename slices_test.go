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

var (
	nan  = math.NaN()
	inf  = math.Inf(1)
	cnan = complex(nan, nan)
	cinf = complex(inf, inf)
)

func TestEqualFloat64s(t *testing.T) {
	const _tol = 4
	cases := []struct {
		Label    string
		In1, In2 []float64
		In3      float64
		Out      bool
	}{
		{"1", []float64{1, 2}, []float64{1, 2, 3}, _tol, false},
		{"2", []float64{0, 0}, []float64{0, 0}, _tol, true},
		{"3", []float64{0, nan}, []float64{0, nan}, _tol, true},
		{"4", []float64{0.9999, 1.1121}, []float64{0.999, 1.112}, _tol, false},
		{"5", []float64{0.9999, 1.1}, []float64{0.99999, 1.1}, _tol, true},
		{"6", []float64{0.99, 1, 0}, []float64{1, 0.99}, _tol, false},
		{"7", []float64{0, 0}, []float64{0, 0.0000001}, _tol, true},
	}
	Test(t, 0, EqualFloat64s, cases)
}

func TestEqualComplex128s(t *testing.T) {
	const _tol = 4
	cases := []struct {
		Label    string
		In1, In2 []complex128
		In3      float64
		Out      bool
	}{
		{"1", []complex128{1 + 1i, 2}, []complex128{1, 2 - 19191i, 3}, _tol, false},
		{"2", []complex128{0, 1i}, []complex128{0, 1i}, _tol, true},
		{"3", []complex128{0, cnan}, []complex128{0, cnan}, _tol, true},
		{"4", []complex128{0.9999i, 1.1121}, []complex128{0.999i, 1.112}, _tol, false},
		{"5", []complex128{0.9999 + 0.007i, 1.1}, []complex128{0.99999 + 0.007i, 1.1}, _tol, true},
		{"6", []complex128{0.99, 1 - 1i}, []complex128{1, 0.99 - 1i}, _tol, false},
	}
	Test(t, 0, EqualComplex128s, cases)
}

func TestAllFloat64s(t *testing.T) {
	cases := []struct {
		Label string
		In1   []float64
		In2   func(float64) bool
		Out   bool
	}{
		{"AllNaN", []float64{1, 2, nan}, math.IsNaN, false},
		{"AllFinite", []float64{1, 2, inf}, func(x float64) bool { return !math.IsInf(x, 0) }, false},
		{"AllZero", []float64{0, 0, 0}, func(x float64) bool { return x == 0 }, true},
	}

	Test(t, 0, AllFloat64s, cases)
}

func TestAllComplex128s(t *testing.T) {
	cases := []struct {
		Label string
		In1   []complex128
		In2   func(complex128) bool
		Out   bool
	}{
		{"AllNaN", []complex128{1, 2, cnan}, cmplx.IsNaN, false},
		{"AllFinite", []complex128{1, 2, cinf}, func(x complex128) bool { return !cmplx.IsInf(x) }, false},
		{"AllZero", []complex128{0, 0, 0}, func(x complex128) bool { return x == 0 }, true},
	}

	Test(t, 0, AllComplex128s, cases)
}
