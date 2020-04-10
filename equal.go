// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import (
	"math"
	"math/cmplx"
	"math/rand"
	"reflect"
	"testing/quick"
	"time"
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
// if they are equal according to reflect.DeepEqual.
//
// The tolerance parameter is optional and has a default value of zero,
// which corresponds to checking for absolute equality.
func Equal(x, y interface{}, tolerance ...interface{}) bool {
	tol := validateTolerance(tolerance...)
	_, ok := equal(reflect.ValueOf(x), reflect.ValueOf(y), tol)
	return ok
}

func equal(xv, yv reflect.Value, tol float64) (i int, ok bool) {
	t := xv.Type()
	if ok = (t.Kind() == yv.Type().Kind()); !ok {
		i--
		return
	}

	switch kx := t.Kind(); kx {
	case reflect.Slice, reflect.Array:
		if ok = equalSlice(xv, yv, tol); !ok {
			return
		}

	case reflect.Map:
		if ok = equalMap(xv, yv, tol); !ok {
			return
		}

	case reflect.Struct:
		if ok = equalStruct(xv, yv, tol); !ok {
			return
		}

	case reflect.Float32, reflect.Float64:
		x := xv.Interface().(float64)
		y := yv.Interface().(float64)
		if ok = equalFloat(x, y, tol); !ok {
			return
		}
	case reflect.Complex64, reflect.Complex128:
		x := xv.Interface().(complex128)
		y := yv.Interface().(complex128)
		if ok = equalComplex(x, y, tol); !ok {
			return
		}
	case reflect.Func:
		if ok = equalFunc(xv, yv, tol); !ok {
			return
		}
	default:
		if ok = reflect.DeepEqual(xv.Interface(), yv.Interface()); !ok {
			return
		}
	}
	return
}

func equalSlice(xv, yv reflect.Value, tol float64) (ok bool) {
	// check the slices have equal lengths
	n := xv.Len()
	if ok = (n == yv.Len()); !ok {
		return
	}
	// check that the items at each position are equal
	for i := 0; i < n; i++ {
		if _, ok = equal(xv.Index(i), yv.Index(i), tol); !ok {
			return
		}
	}
	return
}

func equalMap(xv, yv reflect.Value, tol float64) (ok bool) {
	xkeys := xv.MapKeys()
	ykeys := yv.MapKeys()

	// check that x and y have the same number of keys
	n := len(ykeys)
	if ok = len(xkeys) == n; !ok {
		return
	}

	// check that each key in xkeys is in ykeys. Need to iterate over all xkeys
	// for each ykey since ordering of keys from maps is non-deterministic
	for i := 0; i < n; i++ {
		ykey := ykeys[i]
		for _, xkey := range xkeys {
			if _, ok = equal(xkey, ykey, tol); ok {
				break
			}
		}
		// if ykey was not found, return false
		if !ok {
			return
		}
		// if the items for this key are not equal, return false
		if _, ok = equal(xv.MapIndex(ykey), yv.MapIndex(ykey), tol); !ok {
			return
		}
	}
	return
}

func equalStruct(xv, yv reflect.Value, tol float64) (ok bool) {
	// check that x and y have the same number of fields
	n := xv.Type().NumField()
	if ok = (n == yv.Type().NumField()); !ok {
		return
	}
	// check that the fields at each position are equal
	for i := 0; i < n; i++ {
		if _, ok = equal(xv.Field(i), yv.Field(i), tol); !ok {
			return
		}
	}
	return
}

// equalFloat returns true if x and y are equal within the specified tolerance
func equalFloat(x, y, tol float64) bool {
	if math.IsNaN(x) && math.IsNaN(y) {
		return true
	}
	if x == y || math.IsInf(y, 0) {
		return math.Signbit(x) == math.Signbit(y)
	}
	// check relative error, i.e. |expected - actual| / |actual| < |tol|
	diff := math.Abs(x - y)
	maxdiff := tol * math.Abs(y)
	// if y = 0 or tol*y underflows, set maxdiff = tol
	if maxdiff == 0 {
		maxdiff = tol
	}
	return diff <= maxdiff
}

// equalComplex returns true if x and y are equal within the specified tolerance
// for the real and imaginary parts
func equalComplex(x, y complex128, tol float64) bool {
	return equalFloat(real(x), real(y), tol) && equalFloat(imag(x), imag(y), tol)
}

func equalFunc(xv, yv reflect.Value, tol float64) (ok bool) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	// if checking for exact equality
	if tol == 0 {
		err := quick.CheckEqual(xv.Interface(), yv.Interface(), &quick.Config{Rand: r})
		return err == nil
	}

	for n := 0; n < 1000; n++ {
		args := mockArgs(xv, r)
		xcall := xv.Call(args)
		ycall := yv.Call(args)
		for i := 0; i < len(xcall); i++ {
			if _, ok = equal(xcall[i], ycall[i], tol); !ok {
				return
			}
		}
	}
	return
}

func mockArgs(fv funcv, r *rand.Rand) []reflect.Value {
	nIn := fv.Type().NumIn()
	args := make([]reflect.Value, nIn)
	for i := 0; i < nIn; i++ {
		v, ok := quick.Value(fv.Type().In(i), r)
		panicIf(!ok, "Error. Could not generate mock arguments.")
		args[i] = v
	}
	return args
}

func validateTolerance(tolerance ...interface{}) float64 {
	var tol float64
	if len(tolerance) > 0 {
		switch t := tolerance[0].(type) {
		case float32:
			tol = math.Abs(float64(t))
		case float64:
			tol = math.Abs(t)
		case complex64:
			tol = cmplx.Abs(complex128(t))
		case complex128:
			tol = cmplx.Abs(t)
		case int:
			tol = math.Abs(float64(t))
		case int8:
			tol = math.Abs(float64(t))
		case int16:
			tol = math.Abs(float64(t))
		case int32:
			tol = math.Abs(float64(t))
		case int64:
			tol = math.Abs(float64(t))
		case uint:
			tol = math.Abs(float64(t))
		case uint8:
			tol = math.Abs(float64(t))
		case uint16:
			tol = math.Abs(float64(t))
		case uint32:
			tol = math.Abs(float64(t))
		case uint64:
			tol = math.Abs(float64(t))
		case uintptr:
			tol = math.Abs(float64(t))
		}
	}
	if math.IsInf(tol, 0) || math.IsNaN(tol) {
		tol = 0
	}
	return tol
}
