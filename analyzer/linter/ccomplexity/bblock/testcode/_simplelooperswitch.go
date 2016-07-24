// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {
	// BB #0 ending.

	for i := 0; ; i++ {
		// BB #1 ending.
		switch i { // BB #2 ending.
		case 0:
			fmt.Printf("Zero") // BB #3 ending.
		case 1:
			fmt.Printf("One") // BB #4 ending.
		}
	}
} // BB #5 ending.
