// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

// Unit-test, ignore other errors then what we are testing.
// @SuppressRule("ERROR_IGNORED")
func main() {
	fmt.Print("Printing with fmt.Print()")
	fmt.Printf("Printing with fmt.Printf()")
	fmt.Println("Printing with fmt.Println()")
}
