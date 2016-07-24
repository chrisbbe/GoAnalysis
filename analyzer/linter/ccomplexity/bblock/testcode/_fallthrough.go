// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {

	switch 5 {

	case 0:
		fmt.Println("Zero")
	case 1:
		fmt.Println("One")
		fallthrough
	case 2:
		fmt.Println("Two")
		fallthrough
	case 3:
		fmt.Println("Three")
	case 4:
		fmt.Println("Foure")
	case 5:
		fmt.Println("Five")
		fallthrough
	case 6:
		fmt.Println("Six")
		fallthrough
	default:
		fmt.Println("Default")

	}
}
