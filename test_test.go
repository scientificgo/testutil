// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil_test

import (
	"math"
	"testing"

	. "scientificgo.org/testutil"
)

type f4struct struct {
	Integer   int
	Remainder float64
}

func f1(x float64) float64        { return x * x }
func f2(x complex128) complex128  { return x * x }
func f3(s string) bool            { return len(s) < 6 && len(s) > 3 }
func f4(x float64) f4struct       { return f4struct{int(x), x - float64(int(x))} }
func f5a(x float64) float64       { return 2 * x }
func f5b(x float64) float64       { return x + x }
func f6(n int, x float64) float64 { return math.Jn(n, x) }
func f7(n int, x float64) []float64 {
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = math.Jn(i+1, x)
	}
	return res
}

func TestFunc1(t *testing.T) {
	cases := []struct {
		Label   string
		In, Out interface{}
	}{
		{"1", 0.1, f1(0.1)},
		{"2", 0.2, f1(0.2)},
		{"3", 0.3, f1(0.3)},
	}
	Test(t, acc, cases, f1)
}

func TestFunc2(t *testing.T) {
	cases := []struct {
		Label   string
		In, Out interface{}
	}{
		{"1", 1i, f2(1i)},
		{"2", 0.5 + 0.5i, f2(0.5 + 0.5i)},
		{"3", 1 - 1i, f2(1 - 1i)},
	}
	Test(t, acc, cases, f2)
}

func TestFunc3(t *testing.T) {
	cases := []struct {
		Label   string
		In, Out interface{}
	}{
		{"1", "dog", false},
		{"2", "caterpillar", false},
		{"3", "frog", true},
	}

	Test(t, acc, cases, f3)
}

func TestFunc4(t *testing.T) {
	cases := []struct {
		Label   string
		In, Out interface{}
	}{
		{"1", 1.0, f4struct{1, 0}},
		{"2", 1.1, f4struct{1, 0.1}},
		{"3", math.Pi, f4struct{3, math.Pi - 3}},
	}

	Test(t, acc, cases, f4)
}

func TestFunc5(t *testing.T) {
	cases := []struct {
		Label string
		In    interface{}
	}{
		{"1", 1.0},
		{"2", 1.1},
		{"3", math.Pi},
		{"4", nan},
	}

	Test(t, acc, cases, f5a, f5b)
}

func TestFunc6(t *testing.T) {
	cases := []struct {
		Label         string
		In1, In2, Out interface{}
	}{
		{"1", 1, 0.1, f6(1, 0.1)},
		{"1", 2, 0.2, f6(2, 0.2)},
		{"3", 3, 0.3, f6(3, 0.3)},
	}

	Test(t, acc, cases, f6)
}

func TestFunc7(t *testing.T) {
	cases := []struct {
		Label         string
		In1, In2, Out interface{}
	}{
		{"1", 1, 0.1, f7(1, 0.1)},
		{"1", 2, 0.2, f7(2, 0.2)},
		{"3", 3, 0.3, f7(3, 0.3)},
	}

	Test(t, acc, cases, f7)
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Did not panic as expected.")
		}
	}()
	cases := []struct {
		Label         string
		In1, In2, Out interface{}
	}{
		{"1", 1, 0.1, f1(0.1)},
		{"1", 2, 0.2, f1(0.2)},
		{"3", 3, 0.3, f1(0.3)},
	}

	Test(t, acc, cases, f1)
}
