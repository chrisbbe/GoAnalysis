// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"bytes"
	"log"
	"math/rand"
	"os"
)

func main() {
	if true {
		log.Println("Always true")
	}

	if false {
		log.Println("Should never be printed")
	}

	if 1 >= 2 {
		log.Println("Not possible")
	}

	if 2 >= 1 && 1 >= 0 || 0 >= 0 {
		log.Println("Possible")
	}

	aString := "a"
	bString := "b"

	if aString == bString {
		log.Println("This should never be printed")
	}

	var buf bytes.Buffer
	if buf.Bytes == nil {
		// Not that one is comparing function, not function result, () missed.
		log.Println("This is alway false")
	} else if true {
		log.Println("Might be printed")
	}

	if 1 == rand.Intn(10) {
		log.Println("Result is 1")
	}

	var Line interface{} = "This is a line"
	if value, ok := Line.(string); ok {
		// Should not be flagged as static condition.
		log.Printf("Value contains: %v\n", value)
	}

	result := os.Args[0] == "main.go"
	if result {
		// Should not be flagged as static condition.
		log.Println("Arg[0] = main.go")
	}

}

func GetBiggest(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
