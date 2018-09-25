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

package utils

// Horner returns the polynomial defined by the coefficients c, evaluated at x,
// using Horner's method, i.e.
//
//		      n-1
//	Horner(x, c) = ∑ c[k] x**k = c[0] + c[1]*x + c[2]*x**2 + ... + c[n-1]*x**(n-1)
//		      k=0
//
func Horner(x float64, c ...float64) float64 {
	if c == nil {
		return 0
	}
	n := len(c)
	res := c[n-1]
	for k := n - 2; k >= 0; k-- {
		res = res*x + c[k]
	}
	return res
}
