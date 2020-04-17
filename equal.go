// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"reflect"
	"testing/quick"
	"time"
)

// EqualResult represents the result of an Equal comparison
// between arbitrary x and y.
type EqualResult struct {
	// Ok is true if x equals y.
	Ok bool

	// Numerical is true if x and y are numerical.
	Numerical bool

	// RelativeError is the error in x relative to y if they are numerical.
	// It is a complex number if x and y are complex numbers.
	RelativeError reflect.Value

	// AbsoluteError is the difference between x and y if they are numerical.
	// It is a complex number if x and y are complex numbers.
	AbsoluteError reflect.Value

	// Position is the the first "location" that x does not equal y
	// if x and y are structured data types.
	//
	// For slices and arrays it is the index of first element from x that does not equal y.
	// For structs it is the index of the first field for which x does not equal y.
	// For maps it is the index of the first key for which x does not equal y.
	//
	// If MissingValue is true, Position gives the index in y of the missing field or key.
	Position int

	// LengthMismatch is true if the number of elements, fields or keys in x
	// differs from y for structured data types.
	LengthMismatch bool

	// MissingValue is true if x and y are maps or structs and x is missing one of the keys
	// or fields in y.
	MissingValue bool
}

// Equal reports whether x (actual) is equal to y (expected).
//
// For numerical types, x is equal to y if:
//  |x - y| < tolerance * |y|, for y ≠ 0 (relative error)
//  |x| < tolerance,           for y = 0 (absolute error)
//
// For structured types (slice, array, struct, map), x equals y if
// every element/field/key of x equals that in y.
//
// For func types, x equals y if x(args) equals y(args) for
// randomly generated args.
//
// For other types x equals y if reflect.DeepEqual(x, y) is true.
//
func Equal(x, y, tolerance interface{}) EqualResult {
	tol := validateTolerance(tolerance)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), tol)
}

var floatType = reflect.ValueOf(float64(1)).Type()
var complexType = reflect.ValueOf(complex128(1)).Type()

// equal reports whether the value represented by xv equals that which
// is represented by yv. It recurses through nested structures to compare
// every part for equality. Numerical values are considered equal if they
// are equal within the specified tolerance, which means that x is equal
// to y if and only if
//
//  |x - y| < tol * |y|, for y ≠ 0 (relative error)
//  |x| < tol,           for y = 0 (absolute error)
//
// for floats and for both the real and imaginary parts for complex types.
func equal(xv, yv reflect.Value, tol float64) (res EqualResult) {
	// this occurs when the expected output for y is nil, e.g. for errors,
	// which does not have a concrete type. To avoid panicking, we cast y as
	// a zero of type x. For the example case of errors, this would
	// set y as the zero value for the error type.
	switch {
	case !yv.IsValid() && xv.IsValid():
		yv = reflect.Zero(xv.Type())
	case !xv.IsValid():
		res.Ok = false
		return
	}

	kind := xv.Type().Kind()

	res.RelativeError = reflect.ValueOf(0.)
	res.AbsoluteError = reflect.ValueOf(0.)

	if res.Ok = (kind == yv.Type().Kind()); !res.Ok {
		return
	}

	switch kind {
	case reflect.Slice, reflect.Array:
		if res = equalSlice(xv, yv, tol); !res.Ok {
			return
		}

	case reflect.Map:
		if res = equalMap(xv, yv, tol); !res.Ok {
			return
		}

	case reflect.Struct:
		if res = equalStruct(xv, yv, tol); !res.Ok {
			return
		}

	case reflect.Float32, reflect.Float64, // real-valued
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x := xv.Convert(floatType).Interface().(float64)
		y := yv.Convert(floatType).Interface().(float64)
		if res = equalFloat(x, y, tol); !res.Ok {
			return
		}
	case reflect.Complex64, reflect.Complex128: // complex-valued
		x := xv.Convert(complexType).Interface().(complex128)
		y := yv.Convert(complexType).Interface().(complex128)
		if res = equalComplex(x, y, tol); !res.Ok {
			return
		}

	case reflect.Func:
		if res = equalFunc(xv, yv, tol); !res.Ok {
			return
		}

	default: // anything else: Bool, Chan, String, Interface, Ptr, UnsafePtr
		if res.Ok = reflect.DeepEqual(xv.Interface(), yv.Interface()); !res.Ok {
			return
		}
	}
	return
}

// equalSlice reports whether the slice xv is equal to the slice yv. It checks
// the lengths are equal and the values for each index positiona are equal.
// Numerical values must be equal within the specified tolerance.
func equalSlice(xv, yv reflect.Value, tol float64) (res EqualResult) {
	// check the slices have equal lengths
	n := xv.Len()
	if res.Ok = (n == yv.Len()); !res.Ok {
		res.LengthMismatch = true
		return
	}
	// check that the items at each position are equal
	for i := 0; i < n; i++ {
		if res = equal(xv.Index(i), yv.Index(i), tol); !res.Ok {
			res.Position = i
			return
		}
	}
	return
}

// equalMap reports whether the map xn is equal to the map yv
// for every key, and that they identical keys. Numerical values
// must be equal within the specified tolerance.
func equalMap(xv, yv reflect.Value, tol float64) (res EqualResult) {
	xkeys := xv.MapKeys()
	ykeys := yv.MapKeys()

	// check that x and y have the same number of keys
	n := len(ykeys)
	if res.Ok = len(xkeys) == n; !res.Ok {
		res.LengthMismatch = true
		return
	}

	// check that each key in xkeys is in ykeys. Need to iterate over all xkeys
	// for each ykey since ordering of keys from maps is non-deterministic
	for i := 0; i < n; i++ {
		ykey := ykeys[i]
		for _, xkey := range xkeys {
			if res = equal(xkey, ykey, tol); res.Ok {
				break
			}
		}
		// if ykey was not found, return false
		if !res.Ok {
			res.Position = i
			res.MissingValue = true
			return
		}
		// if the items for this key are not equal, return false
		if res = equal(xv.MapIndex(ykey), yv.MapIndex(ykey), tol); !res.Ok {
			res.Position = i
			return
		}
	}
	return
}

// equalStruct reports whether the struct xn is equal to the struct yv
// for every field, and that they identical fields. Numerical values
// must be equal within the specified tolerance.
func equalStruct(xv, yv reflect.Value, tol float64) (res EqualResult) {
	// check that x and y have the same number of fields
	n := xv.Type().NumField()
	if res.Ok = (n == yv.Type().NumField()); !res.Ok {
		res.LengthMismatch = true
		return
	}
	// check that the fields at each position are equal
	for i := 0; i < n; i++ {
		if res.Ok = xv.Type().Field(i).Name == yv.Type().Field(i).Name; !res.Ok {
			res.MissingValue = true
			res.Position = i
			return
		}
		if res = equal(xv.Field(i), yv.Field(i), tol); !res.Ok {
			res.Position = i
			return
		}
	}
	return
}

// equalFloat reports whether x equals y within the specified tolerance.
// Zeros and Infinities are considered equal if they have the same sign.
// NaNs are always considered equal to other NaNs.
func equalFloat(x, y, tol float64) (res EqualResult) {
	diff := x - y
	res.Numerical = true
	res.AbsoluteError = reflect.ValueOf(diff)

	if math.IsNaN(x) && math.IsNaN(y) {
		res.Ok = true
		res.RelativeError = reflect.ValueOf(0.)
		return
	}

	if x == y || math.IsInf(y, 0) {
		res.Ok = math.Signbit(x) == math.Signbit(y)
		if !res.Ok && math.IsInf(y, 0) {
			res.AbsoluteError = reflect.ValueOf(y)
			res.RelativeError = reflect.ValueOf(y)
		}
		if res.Ok {
			res.AbsoluteError = reflect.ValueOf(0.)
			res.RelativeError = reflect.ValueOf(0.)
		}
		return
	}

	// check magnitude of relative error, i.e. |expected - actual| / |actual| < |tol|
	if y == 0 {
		maxdiff := tol
		res.Ok = math.Abs(diff) < math.Abs(maxdiff)
		res.RelativeError = reflect.ValueOf(diff)
		return
	}

	maxdiff := y * tol
	res.Ok = math.Abs(diff) <= math.Abs(maxdiff)
	res.RelativeError = reflect.ValueOf(diff / y)
	return
}

// equalComplex reports whether x equals y within the specified tolerance
// for both the real and imaginary parts.
func equalComplex(x, y complex128, tol float64) (res EqualResult) {
	rr := equalFloat(real(x), real(y), tol)
	ir := equalFloat(imag(x), imag(y), tol)
	relerr := complex(
		rr.RelativeError.Interface().(float64),
		ir.RelativeError.Interface().(float64),
	)
	abserr := complex(
		rr.AbsoluteError.Interface().(float64),
		ir.AbsoluteError.Interface().(float64),
	)
	res.Numerical = true
	res.Ok = rr.Ok && ir.Ok
	res.RelativeError = reflect.ValueOf(relerr)
	res.AbsoluteError = reflect.ValueOf(abserr)
	return
}

// equalFunc reports whether two functions xv and xy are equivalenet by
// comparing their respective outputs on randomly generated inputs.
// Numerical output values must be equal within the specified tolerance.
func equalFunc(xv, yv reflect.Value, tol float64) (res EqualResult) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	// if checking for exact equality just use the testing/quick package
	if tol == 0 {
		err := quick.CheckEqual(xv.Interface(), yv.Interface(), &quick.Config{Rand: r})
		res.Ok = (err == nil)
		return
	}

	// otherwise generate n random sets of arguments and check the functions
	// agree for each set, returning an error if they do not agree to within the tolerance
	for n := 0; n < 1000; n++ {
		args, err := mockArgs(xv, r)
		if res.Ok = (err == nil); !res.Ok {
			return
		}
		xcall := xv.Call(args)
		ycall := yv.Call(args)
		for i := 0; i < len(xcall); i++ {
			if res = equal(xcall[i], ycall[i], tol); !res.Ok {
				return
			}
		}
	}
	return
}

// mockArgs generates mock arguments for calling an arbitrary function fv
// based on its signature.
func mockArgs(fv reflect.Value, r *rand.Rand) (args []reflect.Value, err error) {
	nIn := fv.Type().NumIn()
	args = make([]reflect.Value, nIn)
	for i := 0; i < nIn; i++ {
		v, ok := quick.Value(fv.Type().In(i), r)
		if !ok {
			err = fmt.Errorf("could not generate mock arguments")
			return
		}
		args[i] = v
	}
	return
}

// validateTolerance ensures the tolerance passed is sensibly valued.
func validateTolerance(tolerance interface{}) (tol float64) {
	t := reflect.ValueOf(tolerance)
	switch kind := t.Kind(); kind {
	case reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		tol = math.Abs(t.Convert(floatType).Interface().(float64))
	case reflect.Complex64, reflect.Complex128:
		tol = cmplx.Abs(t.Convert(complexType).Interface().(complex128))
	}
	if math.IsNaN(tol) {
		tol = 0
	}
	return
}
