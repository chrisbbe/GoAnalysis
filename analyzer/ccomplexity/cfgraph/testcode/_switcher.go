// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {
	number := 3
	fmt.Printf("Number %d is %s\n", number, integerToString(2))
}

func integerToString(number int) string {
	switch number {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	default:
		return "Invalid number"
	}
}
