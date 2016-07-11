// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {
	// BB #0 ending.
	var x interface{}
	x = true
	dataType := ""

	switch x.(type) { // BB #1 ending.

	case nil:
		dataType = "nil" // BB #2 ending.
	case int:
		dataType = "int" // BB #3 ending.
	case bool:
		dataType = "bool" // BB #4 ending.
	case string:
		dataType = "string" // BB #5 ending.
	default:
		dataType = "unknown" // BB #6 ending.
	}
	fmt.Printf("Type is: %s\n", dataType)
} // BB #8 ending.
