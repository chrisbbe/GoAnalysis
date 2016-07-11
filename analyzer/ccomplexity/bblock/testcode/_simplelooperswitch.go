// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {

	for i := 0; ; i++ {
		switch i {
		case 0:
			fmt.Printf("Zero")
		case 1:
			fmt.Printf("One")
		}
	}
}
