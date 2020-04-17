// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil_test

import (
	"math"
	"testing"

	. "scientificgo.org/testutil"
)

var (
	nan  = math.NaN()
	inf  = math.Inf(1)
	cnan = complex(nan, nan)
	cinf = complex(inf, inf)
)

func TestEqual(t *testing.T) {
	type mystruct struct {
		Int    int
		String string
		Float  float64
	}
	type mystruct2 struct {
		Int    int
		String string
		Float  float64
		Extra  string
	}
	type mystruct3 struct {
		Int           int
		String        string
		IncorrectName float64
	}

	fn := func(s string, f64s [][]float64) string {
		return s
	}
	fg := func(s string, f64s [][]float64) string {
		return s
	}

	tol := 1e-4
	cases := []struct {
		Label              string
		In1, In2, In3, Out interface{}
	}{
		{"", []float64{1, 2}, []float64{1, 2, 3}, 0, false},
		{"", []float64{0., nan}, []float64{0., nan}, 0, true},
		{"", []float64{0.9, 1.}, []float64{0.8, 1.}, 0, false},
		{"", []float64{0., 0.00000001}, []float64{0., 0.}, tol, true},
		{"", []complex128{0., 0.999999i}, []complex128{0., 1i}, tol, true},
		{"", []complex128{1. - 1i, 3.3}, []complex128{1. - 1i, 3}, tol, false},
		{"", mystruct{1, "ScientificGo", math.E}, mystruct{1, "ScientificGo", math.E}, nil, true},
		{"", mystruct{1, "ScientificGo", math.E}, mystruct{1, "ScientificGopher", math.E}, nil, false},
		{"", mystruct{1, "Hey", math.E}, mystruct2{1, "Hey", math.E, "extra"}, nil, false},
		{"", mystruct{1, "Hey", math.E}, mystruct3{1, "Hey", math.E}, nil, false},
		{"", map[int]int{0: 1, 1: 10, 2: 100}, map[int]int{0: 1, 1: 10, 2: 100}, tol, true},
		{"", map[int]int{0: 1, 1: 11, 2: 100}, map[int]int{0: 1, 1: 10, 2: 100}, tol, false},
		{"", map[int]int{0: 1, 1: 10, 3: 100}, map[int]int{0: 1, 1: 10, 2: 100}, nil, false},
		{"", map[int]int{0: 1, 1: 10, 2: 100, 3: 1000}, map[int]int{0: 1, 1: 10, 2: 100}, nil, false},
		{"", [2]float64{1, 2}, []float64{1, 2}, 0, false},
		{"", math.Jn, math.Jn, 0, true},
		{"", math.Jn, math.Jn, 1e-9, true},
		{"", math.Jn, math.Yn, 0, false},
		{"", math.Jn, math.Yn, 1e-9, false},
		{"", fn, fg, nil, true},
		{"", +inf, +inf, nil, true},
		{"", +inf, -inf, nan, false},
		{"", -inf, +inf, complex64(0), false},
		{"", +inf, 1., math.Inf(-1), true},
		{"", tol / 100, float64(0), tol, true},
	}

	for _, c := range cases {
		t.Run(c.Label, func(t *testing.T) {
			if res := Equal(c.In1, c.In2, c.In3); res.Ok != c.Out {
				t.Errorf("Error: wanted %v, got %v", c.Out, res)
			}
		})
	}
}
