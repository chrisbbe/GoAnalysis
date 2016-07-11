// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// BB #0 ending.
	stop := make(chan int)

	go func() {
		// BB #1 ending.
		os.Stdin.Read(make([]byte, 1))
		stop <- 1
	}()

	fmt.Println("Started timer, press return to stop watch.")
	tick := time.Tick(time.Second)
	for time := 0; ; time++ {
		// BB #2 ending.
		select { // BB #3 ending.
		case <-tick:
			fmt.Println(time) // BB #5 ending.
		case <-stop:
			fmt.Printf("Watch stopped after %d seconds.\n", time)
			return // BB #5 ending.
		}
	}
} // BB #6 ending.
