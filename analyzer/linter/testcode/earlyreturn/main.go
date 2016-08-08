// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"log"
)

func main() {
	log.Println("Hello World")
	a := "Hei"
	return
	log.Printf("Finally done, a = %s\n", a)
}

func swap(a, b interface{}) (c, d interface{}) {
	c = b
	d = a
	return
}

func getBiggest(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
