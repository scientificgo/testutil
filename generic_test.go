// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils_test

import (
	"math"
	"testing"
	. "scientificgo.org/testutils"
)

type f4struct struct {
	Integer   int
	Remainder float64
}

func func1(x float64) float64       { return x * x }
func func2(x complex128) complex128 { return x * x }
func func3(s string) bool           { return len(s) < 6 && len(s) > 3 }
func func4(x float64) f4struct      { return f4struct{int(x), x - float64(int(x))} }

func jnSlice(n int, x float64) []float64 {
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = math.Jn(i+1, x)
	}
	return res
}

func TestFunc1(t *testing.T) {
	var f func(*testing.T, float64, []string, func(float64) float64, []float64, []float64)
	GenerateTest(&f)

	labels := []string{"1", "2", "3"}
	in := []float64{0.1, 0.2, 0.3}
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = func1(v)
	}

	f(t, tol, labels, func1, in, out)
}

func TestFunc2(t *testing.T) {
	var f func(*testing.T, float64, []string, func(complex128) complex128, []complex128, []complex128)
	GenerateTest(&f)

	labels := []string{"1", "2", "3"}
	in := []complex128{1i, 0.5 + 0.5i, 1}
	out := make([]complex128, len(in))
	for i, v := range in {
		out[i] = func2(v)
	}

	f(t, tol, labels, func2, in, out)
}

func TestFunc3(t *testing.T) {
	var f func(*testing.T, float64, []string, func(string) bool, []string, []bool)
	GenerateTest(&f)

	labels := []string{"1", "2", "3"}
	in := []string{"dog", "caterpillar", "frog"}
	out := make([]bool, len(in))
	for i, v := range in {
		out[i] = func3(v)
	}

	f(t, tol, labels, func3, in, out)
}

func TestFunc4(t *testing.T) {
	var f func(*testing.T, float64, []string, func(float64) f4struct, []float64, []f4struct)
	GenerateTest(&f)

	labels := []string{"1", "2", "3"}
	in := []float64{1.0, 1.1, math.Pi}
	out := make([]f4struct, len(in))
	for i, v := range in {
		out[i] = func4(v)
	}

	f(t, tol, labels, func4, in, out)
}

func TestJn(t *testing.T) {
	var f func(*testing.T, float64, []string, func(int, float64) float64, []int, []float64, []float64)
	GenerateTest(&f)

	labels := []string{"1", "2", "3"}
	in1 := []int{1, 2, 3}
	in2 := []float64{0.1, 0.2, 0.3}
	out := make([]float64, len(in1))
	for i := 0; i < len(in1); i++ {
		out[i] = math.Jn(in1[i], in2[i])
	}

	f(t, tol, labels, math.Jn, in1, in2, out)
}

func TestJnSlice(t *testing.T) {
	var f func(*testing.T, float64, []string, func(int, float64) []float64, []int, []float64, [][]float64)
	GenerateTest(&f)

	labels := []string{"1 output", "2 outputs", "3 outputs"}
	in1 := []int{1, 2, 3}
	in2 := []float64{0.1, 0.2, 0.3}
	res := [][]float64{jnSlice(1, 0.1), jnSlice(2, 0.2), jnSlice(3, 0.3)}

	f(t, tol, labels, jnSlice, in1, in2, res)
}

func TestPanicFuncPos(t *testing.T) {
	var f func(*testing.T, float64, []string, []float64, []float64, func(float64) float64)
	GenerateTest(&f)

	labels := []string{"1", "2", "3"}
	in := []float64{0.1, 0.2, 0.3}
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = func1(v)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Error. Did not panic.")
		}
	}()

	f(t, tol, labels, in, out, func1)
}

func TestPanicWrongIO(t *testing.T) {
	var f func(*testing.T, float64, []string, func(int, float64) float64, []float64, []float64)
	GenerateTest(&f)

	labels := []string{"1", "2", "3"}
	in := []float64{0.1, 0.2, 0.3}
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = func1(v)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Error. Did not panic.")
		}
	}()

	f(t, tol, labels, math.Jn, in, out)
}

func TestPanicIOType(t *testing.T) {
	var f func(*testing.T, float64, []string, func(float64) float64, []float64, float64)
	GenerateTest(&f)

	labels := []string{"1", "2", "3"}
	in := []float64{0.1, 0.2, 0.3}
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = func1(v)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Error. Did not panic.")
		}
	}()

	f(t, tol, labels, func1, in, out[0])
}

func TestPanicIOLengthMismatch(t *testing.T) {
	var f func(*testing.T, float64, []string, func(float64) float64, []float64, []float64)
	GenerateTest(&f)

	labels := []string{"1", "2", "3"}
	in := []float64{0.1, 0.2, 0.3}
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = func1(v)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Error. Did not panic.")
		}
	}()

	f(t, tol, labels, func1, in, out[:len(out)-1])
}
