// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main // BB #0 starting.

func main() {
	t := 2 + 3 // BB #1 starting.
	u := 2 - 3

	if 2 < 3 {
		v := 5 + 5 // BB #2 starting.
	} else if 3 > 2 {
		w := 10 + 3 // BB #3 starting.
		v := w - 4
	}

	x := 2 + 3
	y := 2 - 3
}
