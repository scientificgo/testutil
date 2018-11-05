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

func TestAny(t *testing.T) {
	type mystruct struct {
		Int     int
		String  string
		Float64 float64
	}

	cases := []struct {
		Label    string
		In1, In2 interface{}
		Out      bool
	}{
		{"NaN", math.IsNaN, []float64{1, 2, nan}, true},
		{"Infinite", func(x float64) bool { return !math.IsInf(x, 0) }, []float64{1, 2, inf}, true},
		{"Zero", func(x float64) bool { return x == 0 }, []float64{0, 1, 2}, true},
		{"NaNComplex", cmplx.IsNaN, []complex128{1, 2, cnan}, true},
		{"InfiniteComplex", func(x complex128) bool { return !cmplx.IsInf(x) }, []complex128{1, 2, cinf}, true},
		{"ZeroComplex", func(x complex128) bool { return x == 0 }, []complex128{0, 1, 2}, true},
		{"StructHelloString", func(x mystruct) bool { return x.String == "Hello" }, []mystruct{{10, "Hello", math.Pi}, {100, "Hello!", math.Pi * math.Pi}}, true},
		{"StructPiFloat", func(x mystruct) bool { return x.Float64 == math.Pi }, []mystruct{{1, "Hey", math.E}, {2, "Heey", math.Ln2}}, false},
	}
	Test(t, 0.0, cases, Any)
}

func TestAll(t *testing.T) {
	type mystruct struct {
		Int     int
		String  string
		Float64 float64
	}

	cases := []struct {
		Label    string
		In1, In2 interface{}
		Out      bool
	}{
		{"NaN", math.IsNaN, []float64{1, 2, nan}, false},
		{"Infinite", func(x float64) bool { return !math.IsInf(x, 0) }, []float64{1, 2, inf}, false},
		{"Zero", func(x float64) bool { return x == 0 }, []float64{0, 0, 0}, true},
		{"NaNComplex", cmplx.IsNaN, []complex128{cnan, cnan, cnan}, true},
		{"InfiniteComplex", func(x complex128) bool { return !cmplx.IsInf(x) }, []complex128{1, 2, cinf}, false},
		{"ZeroComplex", func(x complex128) bool { return x == 0 }, []complex128{0, 1, 2}, false},
		{"StructHelloString", func(x mystruct) bool { return x.String == "Hello" }, []mystruct{{10, "Hello", math.Pi}, {100, "Hello", math.Pi * math.Pi}}, true},
		{"StructPiFloat", func(x mystruct) bool { return x.Float64 == math.Pi }, []mystruct{{1, "Hey", math.E}, {2, "Heey", math.Ln2}}, false},
	}
	Test(t, 0.0, cases, All)
}
