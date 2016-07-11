// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

func main() {
	counter := 0

LOOP:
	if 100 > counter {
		counter++
		goto LOOP
	}

}
