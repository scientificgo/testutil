// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import (
	"math"
	"math/rand"
	"reflect"
	"testing/quick"
)

// Equal returns true if x and y are equal.
//
// For float and complex number types, x and y are equal when
// they are equal to the specified number of digits.
//
// For composite types (slice, array, struct, map) x and y are
// equal if they are equal for every element/field/key.
//
// For func types, x and y are equal if their outputs are equal
// when evaluated with mock arguments.
//
// For other types (int, bool, string, char, etc.) x and y are equal
// if x == y is true.
func Equal(x, y interface{}, digits float64) bool {
	_, ok := equal(reflect.ValueOf(x), reflect.ValueOf(y), digits)
	return ok
}

func equal(x, y reflect.Value, digits float64) (i int, ok bool) {
	tx := x.Type()
	if ok = (tx.Kind() == y.Type().Kind()); !ok {
		i--
		goto end
	}

	switch kx := tx.Kind(); kx {
	case reflect.Slice, reflect.Array:
		n := x.Len()
		if ok = (n == y.Len()); !ok {
			i--
			goto end
		}
		for i = 0; i < n; i++ {
			if _, ok = equal(x.Index(i), y.Index(i), digits); !ok {
				goto end
			}
		}
	case reflect.Map:
		keys := x.MapKeys()
		n := len(keys)
		if ok = (len(y.MapKeys()) == n); !ok {
			goto end
		}
		for i = 0; i < n; i++ {
			if _, ok = equal(x.MapIndex(keys[i]), y.MapIndex(keys[i]), digits); !ok {
				goto end
			}
		}
	case reflect.Struct:
		n := tx.NumField()
		if ok = (n == y.Type().NumField()); !ok {
			i--
			goto end
		}
		for i = 0; i < n; i++ {
			if _, ok = equal(x.Field(i), y.Field(i), digits); !ok {
				goto end
			}
		}
	case reflect.Float32, reflect.Float64:
		xf := x.Interface().(float64)
		yf := y.Interface().(float64)
		if ok = equalFloat(xf, yf, digits); !ok {
			goto end
		}
	case reflect.Complex64, reflect.Complex128:
		xc := x.Interface().(complex128)
		yc := y.Interface().(complex128)
		if ok = equalComplex(xc, yc, digits); !ok {
			goto end
		}
	case reflect.Func:
		args := mockArgs(x)
		xcall := x.Call(args)
		ycall := y.Call(args)
		for i := 0; i < len(xcall); i++ {
			if _, ok = equal(xcall[i], ycall[i], digits); !ok {
				goto end
			}
		}
	default:
		if ok = reflect.DeepEqual(x.Interface(), y.Interface()); !ok {
			goto end
		}
	}
end:
	return i, ok
}

// equalFloat returns true if x and y are equal within the specified tolerance
// (i.e. to tol significant figures).
func equalFloat(x, y, digits float64) bool {
	if x == y || equalNaN(x, y) {
		return true
	}
	if x != 0 && y != 0 && sd(1-x/y) >= digits {
		return true
	}
	if sd(x-y) >= digits {
		return true
	}
	return false
}

// equalComplex returns true if x and y are equal to the given number of significant digits.
func equalComplex(x, y complex128, digits float64) bool {
	return equalFloat(real(x), real(y), digits) && equalFloat(imag(x), imag(y), digits)
}

// sd returns the position of the first significant digits of |x|; e.g. sd(1e-17) = 17.
func sd(x float64) float64 {
	if x < 0 {
		x = -x
	}
	return -math.Log10(x)
}

// equalNaN returns true if x and y are NaN.
func equalNaN(x, y float64) bool {
	return math.IsNaN(x) && math.IsNaN(y)
}

func mockArgs(x funcv) []reflect.Value {
	r := rand.New(rand.NewSource(99))
	nIn := x.Type().NumIn()
	args := make([]reflect.Value, nIn)
	for i := 0; i < nIn; i++ {
		v, ok := quick.Value(x.Type().In(i), r)
		panicIf(!ok, "Error. Could not generate mock arguments.")
		args[i] = v
	}
	return args
}
