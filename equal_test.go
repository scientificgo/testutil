package testutils_test

import (
	"math"
	. "scientificgo.org/testutils"
	"testing"
)

var (
	acc  = 10.0
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
	_tol := 4.0
	cases := []struct {
		Label              string
		In1, In2, In3, Out interface{}
	}{
		{"UnequalLengthSlice", []float64{1, 2}, []float64{1, 2, 3}, _tol, false},
		{"", []float64{0, nan}, []float64{0, nan}, _tol, true},
		{"", []float64{0.9999, 1.1121}, []float64{0.999, 1.112}, _tol, false},
		{"", []float64{0, 0}, []float64{0, 0.0000001}, _tol, true},
		{"", []complex128{1 + 1i, 2}, []complex128{1, 2 - 19191i, 3}, _tol, false},
		{"", []complex128{0, 1i}, []complex128{0, 1i}, _tol, true},
		{"", []complex128{0.9999i, 1.1121}, []complex128{0.999i, 1.112}, _tol, false},
		{"", []complex128{0.9999 + 0.007i, 1.1}, []complex128{0.99999 + 0.007i, 1.1}, _tol, true},
		{"", []complex128{0.99, 1 - 1i}, []complex128{1, 0.99 - 1i}, _tol, false},
		{"", mystruct{1, "Hey", math.E}, mystruct{1, "Hey", math.E}, _tol, true},
		{"", mystruct{1, "Hey", math.E}, mystruct{2, "Heey", math.Ln2}, _tol, false},
		{"DifferentStruct", mystruct{1, "Hey", math.E}, mystruct2{1, "Hey", math.E, "extra"}, _tol, false},
		{"", map[int]int{0: 1, 1: 10, 2: 100}, map[int]int{0: 1, 1: 10, 2: 100}, _tol, true},
		{"DifferentMapKeys", map[int]int{0: 1, 1: 10, 2: 100}, map[int]int{0: 1, 1: 10, 2: 100, 3: 1000}, _tol, false},
		{"DifferentMapValues", map[int]int{0: 1, 1: 100, 2: 100}, map[int]int{0: 1, 1: 10, 2: 100}, _tol, false},
		{"DifferentTypes", [2]float64{1, 2}, []float64{1, 2}, _tol, false},
		{"", math.Jn, math.Jn, _tol, true},
		{"", math.Jn, math.Yn, _tol, false},
	}
	Test(t, float64(0), Equal, cases)
}
