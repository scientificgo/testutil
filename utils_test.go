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

package utils_test

import (
	"fmt"
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

func TestHorner(t *testing.T) {
	cases := [][]float64{
		// Each case is: {x, c[0], c[1], ..., out}
		{1, 0},
		{1, 2, 2, 2, 6},
		{2.2011, 1, 2, 3, 4, 5, 6, 489.946813561314360763060},
	}
	for _, cc := range cases {
		n := len(cc) - 2
		x := cc[0]
		out := cc[n+1]
		var c []float64
		if n != 0 {
			c = cc[1 : n+1]
		}
		t.Run(fmt.Sprintf("n=%v", n), func(t *testing.T) {
			res := Horner(x, c...)
			if !EqualFloat64(res, out, tol) {
				t.Errorf("Horner(%v, %v) = %v, want %v", x, c, res, out)
			}
		})
	}
}

func TestReduceFloat64s(t *testing.T) {
	cases := []struct {
		name string
		in   [][]float64
		out  [][]float64
		len  []int
	}{
		{"NaNs", [][]float64{[]float64{NaN, 2, 1, NaN, 5, NaN, 3}, []float64{NaN, NaN, 2, 11, NaN, NaN}},
			[][]float64{[]float64{NaN, 1, NaN, 5, NaN, 3}, []float64{NaN, NaN, 11, NaN, NaN}}, []int{6, 5}},
		{"Infs", [][]float64{[]float64{1, 2, 3, 4, Inf, -Inf}, []float64{-Inf, 11, 22, 33, Inf}},
			[][]float64{[]float64{1, 2, 3, 4}, []float64{11, 22, 33}}, []int{4, 3}},
		{"None", [][]float64{[]float64{1, 2, 3}, []float64{4, 5, 6}},
			[][]float64{[]float64{1, 2, 3}, []float64{4, 5, 6}}, []int{3, 3}},
		{"All", [][]float64{[]float64{math.Pi, math.E, math.Ln2}, []float64{math.Pi, math.E, math.Ln2}},
			[][]float64{[]float64{}, []float64{}}, []int{0, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a, b := c.in[0], c.in[1]
			aout, bout := c.out[0], c.out[1]
			p, q := c.len[0], c.len[1]
			aa, bb, p, q := ReduceFloat64s(a, b)
			if !EqualFloat64s(aa, aout, tol) ||
				!EqualFloat64s(bb, bout, tol) ||
				p != len(aa) || q != len(bb) {
				t.Errorf("ReduceFloat64s(%v, %v) = %v %v %v %v, want %v %v %v %v", a, b, aa, bb, len(aa), len(bb), aout, bout, p, q)
			}
		})
	}
}

func TestReduceComplex128s(t *testing.T) {
	cases := []struct {
		name string
		in   [][]complex128
		out  [][]complex128
		len  []int
	}{
		{"NaNs", [][]complex128{[]complex128{cNaN, 2, 1, cNaN, 5, cNaN, 3}, []complex128{cNaN, cNaN, 2, 11, cNaN, cNaN}},
			[][]complex128{[]complex128{cNaN, 1, cNaN, 5, cNaN, 3}, []complex128{cNaN, cNaN, 11, cNaN, cNaN}}, []int{6, 5}},
		{"Infs", [][]complex128{[]complex128{1, 2, 3, 4, cInf, -cInf}, []complex128{-cInf, 11, 22, 33, cInf}},
			[][]complex128{[]complex128{1, 2, 3, 4}, []complex128{11, 22, 33}}, []int{4, 3}},
		{"None", [][]complex128{[]complex128{1, 2, 3}, []complex128{4, 5, 6}},
			[][]complex128{[]complex128{1, 2, 3}, []complex128{4, 5, 6}}, []int{3, 3}},
		{"All", [][]complex128{[]complex128{math.Pi, math.E, math.Ln2}, []complex128{math.Pi, math.E, math.Ln2}},
			[][]complex128{[]complex128{}, []complex128{}}, []int{0, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a, b := c.in[0], c.in[1]
			aout, bout := c.out[0], c.out[1]
			p, q := c.len[0], c.len[1]
			aa, bb, p, q := ReduceComplex128s(a, b)
			if !EqualComplex128s(aa, aout, tol) ||
				!EqualComplex128s(bb, bout, tol) ||
				p != len(aa) || q != len(bb) {
				t.Errorf("ReduceComplex128s(%v, %v) = %v %v %v %v, want %v %v %v %v", a, b, aa, bb, len(aa), len(bb), aout, bout, p, q)
			}
		})
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
