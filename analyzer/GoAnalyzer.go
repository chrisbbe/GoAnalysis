// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"flag"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var sourceRootDir = flag.String("dir", "", "Absolute path to root directory of Golang source files to be analysed")
var printHelp = flag.Bool("help", false, "Print this usage help")

func main() {
	flag.Parse()

	if flag.NFlag() < 1 {
		flag.Usage()
	}

	// Option dir selected.
	if len(*sourceRootDir) > 0 {
		start := time.Now()

		if ok, err := pathExists(*sourceRootDir); ok {
			goFiles, _ := getGoFiles(*sourceRootDir)
			violations := detectViolations(goFiles)

			log.Print("**** VIOLATIONS ****\n")

			log.Printf("Found %d violations!\n", len(violations))
			for _, vio := range violations {
				log.Printf("Detected violation %32s on line %5d in %s\n", vio.Type, vio.SrcLine, vio.SrcPath)
			}

			elapsed := time.Since(start)
			log.Printf("Analysis took %s\n", elapsed)
		} else {
			log.Print(err)
		}
	}

	// Print help.
	if *printHelp {
		flag.Usage()
	}
}

func detectViolations(goFilePath []string) (violations []*linter.Violation) {
	for _, goFile := range goFilePath {
		violation, err := linter.DetectViolations(goFile)
		if err != nil {
			log.Print(err)
		}

		violations = append(violations, violation...)
	}
	return violations
}

// getGoFiles searches recursively for .go files in the searchDir path, returning the absolute path to the files.
func getGoFiles(searchDir string) (goFiles []string, errors []error) {
	filepath.Walk(searchDir, func(pat string, file os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, err)
		}
		if file != nil && !file.IsDir() && file.Mode().IsRegular() && strings.EqualFold(path.Ext(file.Name()), ".go") {
			// Only add if regular Go source-code file.
			goFiles = append(goFiles, pat)
		}
		return nil
	})
	return goFiles, errors
}

// pathExists returns whether the given file or directory exists or not.
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, err
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return true, err
}
