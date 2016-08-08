// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"errors"
	"log"
)

func main() {
	result, err := divide(2, 2)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Result: %f\n", result)
}

// GitHub issue #4 concerns with false positives detected of ERROR_IGNORED
// on line 23.
func divide(x, y float32) (float32, error) {
	if x == 0.0 || y == 0.0 {
		return -1, errors.New("Zero divison is not smart!")
	}
	return x / y, nil
}
