// Copyright (c) 2015-2016 The GoAnalysis Authors. All rights reserved.
// Use of this source code is governed by the MIT license found in the
// LICENSE file.
package main

import (
	"log"
	"math/rand"
)

func main() {

	if 1 == rand.Intn(10) {
		return
	} else {
		log.Println("Lucky")
	}
}

func foo(e interface{}) interface{} {
	if e == nil {
		return nil
	} else if e != nil {
		return e
	}
	return nil
}
