// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils_test

import (
	"math"
	. "scientificgo.org/testutils"
	"testing"
)

const acc = 10

type f4struct struct {
	Integer   int
	Remainder float64
}

func func1(x float64) float64        { return x * x }
func func2(x complex128) complex128  { return x * x }
func func3(s string) bool            { return len(s) < 6 && len(s) > 3 }
func func4(x float64) f4struct       { return f4struct{int(x), x - float64(int(x))} }
func func5a(x float64) float64       { return 2 * x }
func func5b(x float64) float64       { return x + x }
func func6(n int, x float64) float64 { return math.Jn(n, x) }
func func7(n int, x float64) []float64 {
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = math.Jn(i+1, x)
	}
	return res
}

func TestFunc1(t *testing.T) {
	cases := []struct {
		Label string
		In    float64
		Out   float64
	}{
		{"1", 0.1, func1(0.1)},
		{"2", 0.2, func1(0.2)},
		{"3", 0.3, func1(0.3)},
	}
	Test(t, acc, func1, cases)
}

func TestFunc2(t *testing.T) {
	cases := []struct {
		Label string
		In    complex128
		Out   complex128
	}{
		{"1", 1i, func2(1i)},
		{"2", 0.5 + 0.5i, func2(0.5 + 0.5i)},
		{"3", 1 - 1i, func2(1 - 1i)},
	}
	Test(t, acc, func2, cases)
}

func TestFunc3(t *testing.T) {
	cases := []struct {
		Label string
		In    string
		Out   bool
	}{
		{"1", "dog", false},
		{"2", "caterpillar", false},
		{"3", "frog", true},
	}

	Test(t, acc, func3, cases)
}

func TestFunc4(t *testing.T) {
	cases := []struct {
		Label string
		In    float64
		Out   f4struct
	}{
		{"1", 1.0, f4struct{1, 0}},
		{"2", 1.1, f4struct{1, 0.1}},
		{"3", math.Pi, f4struct{3, math.Pi - 3}},
	}

	Test(t, acc, func4, cases)
}

func TestFunc5(t *testing.T) {
	cases := []struct {
		Label string
		In    float64
	}{
		{"1", 1.0},
		{"2", 1.1},
		{"3", math.Pi},
		{"4", nan},
	}

	Test(t, acc, [2]func(float64) float64{func5a, func5b}, cases)
}

func TestFunc6(t *testing.T) {
	cases := []struct {
		Label string
		In1   int
		In2   float64
		Out   float64
	}{
		{"1", 1, 0.1, func6(1, 0.1)},
		{"1", 2, 0.2, func6(2, 0.2)},
		{"3", 3, 0.3, func6(3, 0.3)},
	}

	Test(t, acc, func6, cases)
}

func TestFunc7a(t *testing.T) {
	cases := []struct {
		Label string
		In1   int
		In2   float64
		Out   []float64
	}{
		{"1", 1, 0.1, func7(1, 0.1)},
		{"1", 2, 0.2, func7(2, 0.2)},
		{"3", 3, 0.3, func7(3, 0.3)},
	}

	Test(t, acc, func7, cases)
}

/* TestFunc7b fails intentionally.
func TestFunc7b(t *testing.T) {
	cases := []struct {
		Label string
		In1   int
		In2   float64
		Out   []float64
	}{
		{"1", 1, 0.1, func7(1, 0.1)},
		{"1", 2, 0.2, func7(2, 0.3)},
		{"3", 3, 0.3, func7(3, 0.2)},
	}

	Test(t, acc, func7, cases)
}
*/

func expectPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("Error. Got nothing, want panic.")
	}
}

func TestPanicCasesNotSlice(t *testing.T) {
	defer expectPanic(t)

	cases := struct {
		Label string
		In    float64
		Out   float64
	}{"1", 0.1, func1(0.1)}
	Test(t, acc, func1, cases)
}
func TestPanicCasesNotStructs(t *testing.T) {
	defer expectPanic(t)

	cases := []float64{1, 2, 3}
	Test(t, acc, func1, cases)
}

func TestPanicNoLabels(t *testing.T) {
	defer expectPanic(t)

	cases := []struct {
		In  float64
		Out float64
	}{
		{0.1, func1(0.1)},
		{0.2, func1(0.2)},
		{0.3, func1(0.3)},
	}
	Test(t, acc, func1, cases)
}

func TestPanicWrongFuncCount(t *testing.T) {
	defer expectPanic(t)

	cases := []struct {
		Label string
		In    float64
		Out   float64
	}{
		{"1", 0.1, func1(0.1)},
		{"2", 0.2, func1(0.2)},
		{"3", 0.3, func1(0.3)},
	}
	Test(t, acc, [3]func(float64) float64{func1, func1, func1}, cases)
}

func TestPanicWrongFuncType(t *testing.T) {
	defer expectPanic(t)

	cases := []struct {
		Label string
		In    float64
		Out   float64
	}{
		{"1", 0.1, func1(0.1)},
		{"2", 0.2, func1(0.2)},
		{"3", 0.3, func1(0.3)},
	}
	Test(t, acc, "func", cases)
}

func TestPanicWrongIOCountOneFunction(t *testing.T) {
	defer expectPanic(t)

	cases := []struct {
		Label string
		In    float64
	}{
		{"1", 0.1},
		{"2", 0.2},
		{"3", 0.3},
	}
	Test(t, acc, func1, cases)
}
func TestPanicWrongIOCountTwoFunctions(t *testing.T) {
	defer expectPanic(t)

	cases := []struct {
		Label string
		In    int
	}{
		{"1", 1},
		{"2", 2},
		{"3", 3},
	}
	Test(t, acc, [2]func(int, float64) float64{func6, func6}, cases)
}
