// Copyright (c) 2018, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutils

import "reflect"

func sliced(c reflect.Value, start, n int) []reflect.Value {
	v := make([]reflect.Value, n)
	for i := 0; i < n; i++ {
		v[i] = underlyingValue(c.Field(start + i))
	}
	return v
}

func underlyingValue(v reflect.Value) reflect.Value {
	if k := v.Kind(); k == reflect.Interface || k == reflect.Ptr {
		return v.Elem()
	}
	return v
}
