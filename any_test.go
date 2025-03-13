// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil_test

import (
	"math"
	"math/cmplx"
	"testing"

	. "github.com/scientificgo/testutil"
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
		Out1     bool
		Out2     error
	}{
		{"",
			math.IsNaN,
			[]float64{1, 2, nan},
			true, nil,
		},

		{"",
			func(x float64) bool { return !math.IsInf(x, 0) },
			[]float64{1, 2, inf},
			true, nil,
		},

		{"",
			func(x float64) bool { return x == 0 },
			[]float64{0, 1, 2},
			true, nil,
		},

		{"",
			cmplx.IsNaN,
			[]complex128{1, 2, cnan},
			true, nil,
		},

		{"",
			func(x complex128) bool { return !cmplx.IsInf(x) },
			[]complex128{1, 2, cinf},
			true, nil,
		},

		{"",
			func(x complex128) bool { return x == 0 },
			[]complex128{0, 1, 2},
			true, nil,
		},

		{"",
			func(x mystruct) bool { return x.String == "Hello" },
			[]mystruct{{10, "Hello", math.Pi}, {100, "Hello!", math.Pi * math.Pi}},
			true, nil,
		},

		{"",
			func(x mystruct) bool { return x.Float64 == math.Pi },
			[]mystruct{{1, "Hey", math.E}, {2, "Heey", math.Ln2}},
			false, nil,
		},
	}
	Test(t, nil, cases, Any)
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
		Out1     bool
		Out2     error
	}{
		{"",
			math.IsNaN,
			[]float64{1, 2, nan},
			false, nil,
		},

		{"",
			func(x float64) bool { return !math.IsInf(x, 0) },
			[]float64{1, 2, inf},
			false, nil,
		},

		{"",
			func(x float64) bool { return x == 0 },
			[]float64{0, 0, 0},
			true, nil,
		},

		{"",
			cmplx.IsNaN,
			[]complex128{cnan, cnan, cnan},
			true, nil,
		},

		{"",
			func(x complex128) bool { return !cmplx.IsInf(x) },
			[]complex128{1, 2, cinf},
			false, nil,
		},

		{"",
			func(x complex128) bool { return x == 0 },
			[]complex128{0, 1, 2},
			false, nil,
		},

		{"",
			func(x mystruct) bool { return x.String == "Hello" },
			[]mystruct{{10, "Hello", math.Pi}, {100, "Hello", math.Pi * math.Pi}},
			true, nil,
		},

		{"",
			func(x mystruct) bool { return x.Float64 == math.Pi },
			[]mystruct{{1, "Hey", math.E}, {2, "Heey", math.Ln2}},
			false, nil,
		},
	}
	Test(t, nil, cases, All)
}
