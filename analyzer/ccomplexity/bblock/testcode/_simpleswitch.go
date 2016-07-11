// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {
	// BB #0 ending.
	number := 3

	switch number { // BB #1 ending.

	case 0: // BB #2 ending.
		fmt.Println("0")
	case 1: // BB #3 ending.
		fmt.Println("1")
	case 2: // BB #4 ending.
		fmt.Println("2")
	default: // BB #5 ending.
		fmt.Println("Invalid number") // BB #6 ending.
	}
}
