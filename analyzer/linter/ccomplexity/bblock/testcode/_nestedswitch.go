// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

const TWO = 2

// @SuppressRule("FMT_PRINTING")
// @SuppressRule("ERROR_IGNORED")
func main() {
	// BB #0 ending
	number := 1
	secondNumber := 3

	switch number { // BB #1 ending.

	case 0:
		fmt.Println("0") // BB #2 ending.
	case 1:
		fmt.Println("2")
		switch secondNumber { // BB #3 ending.

		case 1:
			fmt.Println("1") // BB #4 ending.
		case 2:
			fmt.Println(TWO) // BB #5 ending.
		default:
			fmt.Printf("No match, secondNumber is %d!\n", number) // BB #6 ending.
		}
	case 3:
		fmt.Println("3") // BB #7 ending.
	case 4:
		fmt.Println("4") // BB #8 ending.
	default:
		fmt.Printf("No match, number is %d!\n", number) // BB #9 ending.
	}
} // BB #10 ending.
