// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "log"

func main() {
	line := "FooBar"

	if len(line) == 0 {
		log.Println("line is empty")
	} else {
		log.Println("line is not empty")
	}
}
