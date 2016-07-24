// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

func main() {
	// BB #0 ending.
	t := 2 + 3
	u := 2 - 3

	if 2 < 3 {
		// BB #1 ending.
		v := 5 + 5
	} else {
		// BB #2 ending.
		w := 10 + 3
		v := w - 4
	} // BB #3 ending.
} // BB #4 ending.
