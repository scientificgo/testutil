// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil_test

import (
	"fmt"
	"testing"

	. "github.com/scientificgo/testutil"
)

func TestParseCases(t *testing.T) {
	// good - should parse correctly
	casesHypot := []struct {
		Label         string
		In1, In2, Out float64
	}{
		{"Hypot2", 0., 0., 0.},
		{"Hypot2", 1., 0., 1.},
		{"Hypot2", 3., 4., 5.},
		{"Hypot2", 4., 3., 5.},
	}
	casesAtoi := []struct {
		Label string
		In    string
		Out1  int
		Out2  error
	}{
		{"Atoi", "10000", 10000, nil},
		{"Atoi", "-99", -99, nil},
	}
	casesItoa := []struct {
		Label string
		In    int
		Out   string
	}{
		{"Itoa", 1, "1"},
	}

	// bad - should raise an error
	casesNotSlice := struct {
		Label   string
		In, Out int
	}{
		"", 1, 2,
	}
	casesNone := []struct {
		Label   string
		In, Out int
	}{}
	casesNotStructs := []map[string]interface{}{
		{"Label": "1", "In": 7, "Out": 3},
	}
	casesEmpty := []struct{}{
		{},
		{},
	}
	casesNoLabel := []struct {
		In, Out float64
	}{
		{1., 2.},
	}

	cases := []struct {
		Label           string
		Cases           Cases
		Ncases, Nfields int
		Err             error
	}{
		{"Good", casesHypot, 4, 4, nil},
		{"Good", casesAtoi, 2, 4, nil},
		{"Good", casesItoa, 1, 3, nil},

		{"Bad", casesNotSlice, 0, 0, fmt.Errorf("wrong kind of argument. Got struct, want slice")},
		{"Bad", casesNone, 0, 0, fmt.Errorf("too few cases. Got 0, want at least 1")},
		{"Bad", casesNotStructs, 1, 0, fmt.Errorf("wrong input type. Got []map, want []struct")},
		{"Bad", casesEmpty, 2, 0, fmt.Errorf("too few fields in cases. Got 0, want at least 1")},
		{"Bad", casesNoLabel, 1, 2, fmt.Errorf("invalid type for first field. Got float64, want string")},
	}

	for _, c := range cases {
		t.Run(c.Label, func(t *testing.T) {
			_, nc, nf, err := ParseCases(c.Cases)
			if res := Equal(&err, &c.Err, nil); !res.Ok || nc != c.Ncases || nf != c.Nfields {
				t.Error(err, nc, c.Ncases, nf, c.Nfields)
			}
		})
	}
}
