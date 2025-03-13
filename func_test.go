// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil_test

import (
	"fmt"
	"math"
	"testing"

	. "github.com/scientificgo/testutil"
)

func TestParseFuncs(t *testing.T) {
	cases := []struct {
		Label string
		Funcs []Func
		Err   error
	}{
		{"Good", []Func{math.Abs}, nil},
		{"Good", []Func{math.Sin, math.Cos}, nil},
		{"Bad", []Func{nil}, fmt.Errorf("invalid functions, reflection failed")},
		{"Bad", []Func{"this is not a func"}, fmt.Errorf("wrong kind of argument. Got string, want func")},
		{"Bad", []Func{}, fmt.Errorf("wrong number of functions. Got 0, want 1 or 2")},
		{"Bad", []Func{math.Sin, math.Cos, math.Sincos}, fmt.Errorf("wrong number of functions. Got 3, want 1 or 2")},
	}

	for _, c := range cases {
		t.Run(c.Label, func(t *testing.T) {
			_, _, err := ParseFuncs(c.Funcs...)
			if res := Equal(&c.Err, &err, nil); !res.Ok {
				t.Error(err)
			}
		})
	}
}
