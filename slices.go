// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import "math"

// AnyFloat64s returns true if testf evaluates to true for any element of slice x.
func AnyFloat64s(x []float64, testf func(float64) bool) bool {
	for _, val := range x {
		if testf(val) {
			return true
		}
	}
	return false
}

// AllFloat64s returns true if testf evaluates to true for all elements of slice x.
func AllFloat64s(x []float64, testf func(float64) bool) bool {
	return !AnyFloat64s(x, func(val float64) bool { return !testf(val) })
}

// AnyComplex128s returns true if testf evaluates to true for any element of slice x.
func AnyComplex128s(x []complex128, testf func(complex128) bool) bool {
	for _, val := range x {
		if testf(val) {
			return true
		}
	}
	return false
}

// AllComplex128s returns true if testf evaluates to true for all elements of slice x.
func AllComplex128s(x []complex128, testf func(complex128) bool) bool {
	return !AnyComplex128s(x, func(val complex128) bool { return !testf(val) })
}

// EqualFloat64 returns true if x and y are equal within the specified tolerance
// (i.e. to tol significant figures).
func EqualFloat64(x, y, tol float64) bool {
	if x == y || equalNaN(x, y) {
		return true
	}
	if x != 0 && y != 0 {
		if sd(1-x/y) >= tol {
			return true
		}
	}
	if sd(x-y) >= tol {
		return true
	}
	return false
}

// EqualComplex128 returns true if x and y are equal within the specified tolerance
// (i.e. to tol significant figures).
func EqualComplex128(x, y complex128, tol float64) bool {
	return EqualFloat64(real(x), real(y), tol) && EqualFloat64(imag(x), imag(y), tol)
}

// EqualFloat64s returns true if x and y are equal within the specified tolerance
// (i.e. to tol significant figures).
func EqualFloat64s(x, y []float64, tol float64) bool {
	if len(x) != len(y) {
		return false
	}
	for i := 0; i < len(x); i++ {
		if !EqualFloat64(x[i], y[i], tol) {
			return false
		}
	}
	return true
}

// EqualComplex128s returns true if x and y are equal within the specified tolerance
// (i.e. to tol significant figures).
func EqualComplex128s(x, y []complex128, tol float64) bool {
	if len(x) != len(y) {
		return false
	}
	for i := 0; i < len(x); i++ {
		if !EqualComplex128(x[i], y[i], tol) {
			return false
		}
	}
	return true
}

// sd returns the position of the first significant digits of |x|; e.g. sd(1e-17) = 17.0
func sd(x float64) float64 {
	if x < 0 {
		x = -x
	}
	return -math.Log10(x)
}

// equalNaN returns false if only x or y is NaN.
func equalNaN(x, y float64) bool {
	return math.IsNaN(x) && math.IsNaN(y)
}
