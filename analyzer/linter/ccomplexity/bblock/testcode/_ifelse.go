// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "log"

func main() {
	// BB #0 ending.
	t := 2 + 3
	u := 2 - 3

	if 2 < 3 {
		// BB #1 ending.
		t = 5 + 5
	} else {
		// BB #2 ending.
		u = 10 + 3
		t = u - 4
	} // BB #3 ending.

	log.Printf("t = %d\n", t)
	log.Printf("u = %d\n", u)
} // BB #4 ending.
