// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main // BB #0 starting.
import "log"

func main() {
	t := 2 + 3 // BB #1 starting.
	u := 2 - 3

	if 2 < 3 {
		t = 5 + 5 // BB #2 starting.
	} else if 3 > 2 {
		t = 10 + 3 // BB #3 starting.
		u = t - 4
	}

	t = 2 + 3
	u = 2 - 3

	log.Printf("t = %d\n", t)
	log.Printf("u = %d\n", u)
}
