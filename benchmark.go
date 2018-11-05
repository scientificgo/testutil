// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import "testing"

// Benchmark runs a sub-benchmark for each case in cs using the function(s) in fs.
func Benchmark(b *testing.B, cs Cases, f Func) {
	cvs, nc, nfc := parseCases(cs)
	fv, _ := parseFuncs(f)

	nIn := fv.Type().NumIn()
	panicIf(nfc-1 != nIn, "Wrong number of input slices. Got %v, want %v.", nfc-1, nIn)

	for i := 0; i < nc; i++ {
		subbench(b, cvs.Index(i), fv, nIn)
	}
}

// subbench runs a sub-benchmark for the case cv using function fv.
func subbench(b *testing.B, cv casev, fv funcv, nIn int) {
	inputs := sliceFrom(cv, 1, nIn)
	b.Run(name(cv), func(b *testing.B) {
		for k := 0; k < b.N; k++ {
			_ = fv.Call(inputs)
		}
	})

}
