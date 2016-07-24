// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {
	// BB #0 ending.
	number := 3
	secondNumber := 3

	switch number { // BB #1 ending.

	case 0: // BB #2 ending.
		fmt.Println("0")
	case 1: // BB #4 ending.
		fmt.Println("2")
		switch secondNumber { // BB #5 ending.

		case 1: // BB #6 ending.
			fmt.Println("1")
		case 2: // BB #7 ending.
			fmt.Println("2")
		default: // BB #8 ending.
			fmt.Printf("No match, secondNumber is %d!\n", number)
		}
	case 3: // BB #9 ending.
		fmt.Println("3")

		var x interface{}
		x = true

		switch t := x.(type) { // BB #10 ending.

		case int: // BB #11 ending.
			fmt.Println("Type is int")
		case string: // BB #12 ending.
			fmt.Println("Type is string")
		case bool: // BB #13 ending.
			fmt.Println("Type is bool")
		default: // BB #14 ending.
			fmt.Printf("%t is unknown type.\n", t)
		}

	case 4: // BB #15 ending.
		fmt.Println("4")
	default: // BB #16 ending.
		fmt.Printf("No match, number is %d!\n", number)
	}
}
