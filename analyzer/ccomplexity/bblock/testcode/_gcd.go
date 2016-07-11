// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {
	// BB #0 ending.
	fmt.Println(gcd(33, 77))
	fmt.Println(gcd(49865, 69811))
}// BB #1 ending.

func gcd(x, y int) int {
	// BB #2 ending.
	for y != 0 {
		// BB #3 ending.
		x, y = y, x % y // BB #4 ending.
	}
	return x // BB #5 ending.
}
