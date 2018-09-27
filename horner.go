// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package utils

// Horner returns the polynomial defined by the coefficients c, evaluated at x,
// using Horner's method, i.e.
//
//                n-1
//  Horner(x, c) = ∑ c[k] x**k = c[0] + c[1]*x + c[2]*x**2 + ... + c[n-1]*x**(n-1)
//                k=0
//
func Horner(x float64, c ...float64) float64 {
	if c == nil {
		return 0
	}

	n = len(c)
	res := c[n-1]
	for k := n - 2; k >= 0; k-- {
		res = res*x + c[k]
	}
	return res
}
