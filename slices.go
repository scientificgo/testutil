/*
   SciGo is a scientific library for the Go language.
   Copyright (C) 2018, Jack Parkinson

   This program is free software: you can redistribute it and/or modify it
   under the terms of the GNU Lesser General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package utils

import "math"

// ReduceFloat64 removes values common to both a and b and returns the reduced slices and their lengths.
func ReduceFloat64s(a, b []float64) ([]float64, []float64, int, int) {
	na := len(a)
	nb := len(b)
	for i := 0; i < na; i++ {
		for j := 0; j < nb; j++ {
			if a[i] == b[j] {
				a = append(a[:i], a[i+1:]...)
				b = append(b[:j], b[j+1:]...)
				i--
				j--
				na--
				nb--
				break
			}
		}
	}
	return a, b, na, nb
}

// ReduceComplex128 removes values common to both a and b and returns the reduced slices and their lengths.
func ReduceComplex128s(a, b []complex128) ([]complex128, []complex128, int, int) {
	na := len(a)
	nb := len(b)
	for i := 0; i < na; i++ {
		for j := 0; j < nb; j++ {
			if a[i] == b[j] {
				a = append(a[:i], a[i+1:]...)
				b = append(b[:j], b[j+1:]...)
				i--
				j--
				na--
				nb--
				break
			}
		}
	}
	return a, b, na, nb
}

// Any returns true if testf evaluates to true for any element of slice x.
func AnyFloat64s(x []float64, testf func(float64) bool) bool {
	for _, val := range x {
		if testf(val) {
			return true
		}
	}
	return false
}

// All returns true if testf evaluates to true for all elements of slice x.
func AllFloat64s(x []float64, testf func(float64) bool) bool {
	return !AnyFloat64s(x, func(val float64) bool {
		return !testf(val)
	})
}

// AnyComplex128 returns true if testf evaluates to true for any element of slice x.
func AnyComplex128s(x []complex128, testf func(complex128) bool) bool {
	for _, val := range x {
		if testf(val) {
			return true
		}
	}
	return false
}

// AllComplex128 returns true if testf evaluates to true for all elements of slice x.
func AllComplex128s(x []complex128, testf func(complex128) bool) bool {
	return !AnyComplex128s(x, func(val complex128) bool {
		return !testf(val)
	})
}

// EqualFloat64 returns true if x and y are equal within the specified tolerance.
func EqualFloat64(x, y, tol float64) bool {
	if x == y {
		return true
	}
	if sd(x, y) >= tol {
		return true
	}
	if equalNaN(x, y) {
		return true
	}
	return false
}

// EqualComplex128 returns true if x and y are equal within the specified tolerance.
func EqualComplex128(x, y complex128, tol float64) bool {
	return EqualFloat64(real(x), real(y), tol) && EqualFloat64(imag(x), imag(y), tol)
}

// EqualFloat64s returns true if the slices x and y are equal within the specified tolerance.
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

// EqualComplex128s returns true if x and y are equal within the specified tolerance.
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

// sd returns the number of significant digits to which x and y are equal.
func sd(x, y float64) float64 {
	switch {
	case y == 0 || x == 0:
		v := x
		if x == 0 {
			v = y
		}
		return -math.Log10(math.Abs(v))
	default:
		return -math.Log10(math.Abs(1 - x/y))
	}
}

// equalNaN returns false if only x or y is NaN.
func equalNaN(x, y float64) bool {
	return math.IsNaN(x) && math.IsNaN(y)
}
