// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import (
	"math"
	"math/rand"
	"reflect"
	"testing/quick"
)

// Equal returns true if x (actual) is equal to y (expected).
//
// For float and complex number types, x is equal to y if is
// within tol of y; in practice, this means |x-y| < |y * tol|
// or |x| < |tol| if y=0.
//
// For composite types (slice, array, struct, map), x equals y if
// every element/field/key of x equals that in y.
//
// For func types, x equals y if x(args) is equal to y(args) for
// randomly generated mock args.
//
// For other types (int, bool, string, char, etc.) x equals y
// if they are deeply equal.
//
// If unspecified, tol=1e-10 will be used.
func Equal(x, y interface{}, tol ...float64) bool {
	t := 1.e-10
	if len(tol) > 0 {
		t = tolerance(tol[0])
	}
	_, ok := equal(reflect.ValueOf(x), reflect.ValueOf(y), t)
	return ok
}

func tolerance(tol float64) float64 {
    if math.IsInf(tol, 0) || math.IsNaN(tol) {
        tol = 1.e-10
    }
    if tol < 0 {
        tol = -tol
    }
    return tol
}

func equal(x, y reflect.Value, tol float64) (i int, ok bool) {
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
			if _, ok = equal(x.Index(i), y.Index(i), tol); !ok {
				goto end
			}
		}
	case reflect.Map:
		xkeys := x.MapKeys()
        ykeys := y.MapKeys()
		n := len(ykeys)
		if ok = len(xkeys) == n; !ok {
			goto end
		}
		for i =0; i < n; i++ {
            // check that each key in xkeys is in ykeys. Need to iterate over all xkeys
            // for each ykey since ordering of keys from maps is not deterministic, so
            // the keys could come in different orders even if xkeys and ykeys contain
            // the same values.
            ykey := ykeys[i]
            for _, xkey := range xkeys {
                if _, ok = equal(xkey, ykey, tol); ok {
                    break
                }
            }
            if !ok {
                goto end
            }
			if _, ok = equal(x.MapIndex(ykey), y.MapIndex(ykey), tol); !ok {
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
			if _, ok = equal(x.Field(i), y.Field(i), tol); !ok {
				goto end
			}
		}
	case reflect.Float32, reflect.Float64:
		xf := x.Interface().(float64)
		yf := y.Interface().(float64)
		if ok = equalFloat(xf, yf, tol); !ok {
			goto end
		}
	case reflect.Complex64, reflect.Complex128:
		xc := x.Interface().(complex128)
		yc := y.Interface().(complex128)
		if ok = equalComplex(xc, yc, tol); !ok {
			goto end
		}
	case reflect.Func:
		args := mockArgs(x)
		xcall := x.Call(args)
		ycall := y.Call(args)
		for i := 0; i < len(xcall); i++ {
			if _, ok = equal(xcall[i], ycall[i], tol); !ok {
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
func equalFloat(x, y, tol float64) bool {
	if equalNaN(x, y) {
	    return true
    }
    if x == y {
        return math.Signbit(x) == math.Signbit(y)
    }
    if math.IsInf(y, 0) {
        return x == y 
    }
    diff := math.Abs(x - y)
    err := tol * math.Abs(y)
    // If y = 0, set err = tol.
    if y == 0 {
        err = tol
    }
    return diff <= err
}

// equalComplex returns true if x and y are equal to the given number of significant tol.
func equalComplex(x, y complex128, tol float64) bool {
	return equalFloat(real(x), real(y), tol) && equalFloat(imag(x), imag(y), tol)
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
