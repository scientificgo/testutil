// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils_test

import (
	. "scientificgo.org/testutils"
	"testing"
)

func BenchmarkAny(b *testing.B) {
	cases := []struct {
		Label    string
		In1, In2 interface{}
	}{
		{"float64", func(x float64) bool { return x == 0 }, []float64{0.1, -0.0001, 0.01, 0, 0.000000000000000000001}},
		{"map", func(x map[int]int) bool { return x == nil }, []map[int]int{{}, {0: 1}}},
	}
	Benchmark(b, cases, Any)
}

func BenchmarkEqual(b *testing.B) {
	type mystruct struct {
		A, B int
		C    string
	}
	cases := []struct {
		Label    string
		In1, In2 interface{}
		Tol      float64
	}{
		{"float64", []float64{1e-5, 1e-11, 1e-18, 1e-20}, []float64{1e-5, 1e-11, 1.1e-18, 2e-20}, 10},
		{"[][]float64", [][][]float64{{{1e-5, 1e-11, 1e-18, 1e-20}}}, [][][]float64{{{1e-5, 1e-11, 1.1e-18, 2e-20}}}, 10},
		{"[]struct", [][]mystruct{{{1, 2, "A"}, {3, 4, "B"}}}, [][]mystruct{{{1, 2, "A"}, {3, 4, "B"}}}, 10},
	}
	Benchmark(b, cases, Equal)
}
