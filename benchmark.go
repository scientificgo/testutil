// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import "testing"

// Benchmark is a generic testing tool that benchmarks the outputs of an
// arbitrary function f to the specified number of significant digits
// (for float or complex types).
//
// cases must be a slice of structs that:
// (i) have a first field that is a string;
// (ii) export all fields.
func Benchmark(b *testing.B, f interface{}, cases interface{}) {
	casesv, f1, nCases, nIn := validateBenchmark(f, cases)
	for i := 0; i < nCases; i++ {
		c := casesv.Index(i)
		inputs := sliced(c, 1, nIn)
		b.Run(c.Field(0).String(), func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				_ = f1.Call(inputs)
			}
		})
	}
}
