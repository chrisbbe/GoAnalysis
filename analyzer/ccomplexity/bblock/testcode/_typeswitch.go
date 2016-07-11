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

	case nil: // BB #2 ending.
		dataType = "nil"
	case int: // BB #3 ending.
		dataType = "int"
	case bool: // BB #4 ending.
		dataType = "bool"
	case string: // BB #5 ending.
		dataType = "string"
	default: // BB #6 ending.
		dataType = "unknown"
	}
	fmt.Printf("Type is: %s\n", dataType) // BB #8 ending.
}
