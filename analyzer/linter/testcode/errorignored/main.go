// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// Returning error should be checked with error != nil
	file, _ := os.Open("file.go")
	defer file.Close()

	logger.Println(file)

	// Returning error should not be ignored with _
	file2, _ := OpenFile("log.txt")

	logger.Println(file2)

	// Error returned here should be handled!
	writeToConsole("Hello console")

	err := writeToConsole("Hello console")
	if err != nil {
		panic(err)
	}

	// Should also detect the error not detected here!
	os.Open("file.go")
}

func OpenFile(filePath string) (*os.File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func writeToConsole(line string) error {
	w := bufio.NewWriter(os.Stdout)
	if _, err := w.WriteString(line); err != nil {
		return err
	}
	w.Flush()
	return nil
}
