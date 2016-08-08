// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"log"
)

// This small piece of code illustrates
// one of the common mistakes new Go
// people do when using Goroutines inside
// for-loops, referencing the loop-variable.
//
func main() {
	//Not-thread safe.
	for num := 0; num < 5; num++ {
		go func() {
			log.Printf("Goroutine #%d\n", num)
		}()
	}

	//Safe loop, no Go routine!
	for val := 0; val < 50; val++ {
		log.Println(val)
	}

	//Thread safe loop.
	for val := 0; val < 5; val++ {
		go func(val int) {
			log.Println(val)
		}(val)
	}

}
