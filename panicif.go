// Copyright (c) 2020, Jack Parkinson. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package testutil

import "fmt"

// panicIf panics if b is true using the format string s and values v.
func panicIf(b bool, s string, v ...interface{}) {
	if b {
		panic(fmt.Sprintf(s, v...))
	}
}
