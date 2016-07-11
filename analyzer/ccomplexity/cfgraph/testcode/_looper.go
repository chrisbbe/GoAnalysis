// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {
	// BB #0 ending.
	sum := 0
	for i := 0; i < 10; i++ {
		// BB #1 ending.
		sum += i // BB #2 ending.
	}
	fmt.Println(sum) // BB #3 ending.
}
