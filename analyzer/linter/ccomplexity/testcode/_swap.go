// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {
	stringA := "A"
	stringB := "B"

	fmt.Printf("Before swap:\t stringA = %s, stringB = %s\n", stringA, stringB)
	swap(&stringA, &stringB)
	fmt.Printf("After swap:\t stringA = %s, stringB = %s\n", stringA, stringB)
}

func swap(a *string, b *string) {
	tmp := a
	*a = *b
	*b = *tmp
}
