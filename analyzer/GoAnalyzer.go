// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"encoding/json"
	"flag"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var sourceRootDir = flag.String("dir", "", "Absolute path to root directory of Golang source files to be analysed.")
var jsonOutput = flag.Bool("json", false, "Print result as JSON.")
var printHelp = flag.Bool("help", false, "Print this usage help.")

func main() {
	flag.Parse()

	if flag.NFlag() < 1 {
		flag.Usage()
	}

	//var violations []*linter.GoFile
	// Option dir selected.
	if len(*sourceRootDir) > 0 {
		start := time.Now()

		if ok, err := pathExists(*sourceRootDir); ok {
			goFiles, errors := getGoFiles(*sourceRootDir)

			for _, err := range errors {
				if err != nil {
					log.Fatal(err)
				}
			}

			goFileViolations, err := linter.DetectViolations(goFiles...)
			if err != nil {
				log.Fatal(err)
			}

			// Direct output to console as JSON.
			if *jsonOutput && len(goFileViolations) > 0 {
				json, err := json.MarshalIndent(goFileViolations, "", "\t")
				if err != nil {
					log.Fatal(err)
				}
				os.Stdout.Write(json)
			} else {
				log.SetOutput(os.Stdout) // We want to send output to stdout, instead of Stderr.
				numberOfViolations := 0
				// Nicely print the output to the console.
				log.Println("-----------------------------------------------------------------------------------------------")

				for _, goFile := range goFileViolations {
					log.Printf("Violations in %s :\n", goFile.FilePath)
					numberOfViolations += len(goFile.Violations)

					for i, violation := range goFile.Violations {
						log.Printf("\t%d) %s (Line %d) - %s\n", i, violation.Type, violation.SrcLine, violation.Description)
					}
					log.Println("-----------------------------------------------------------------------------------------------")
				}
				log.Printf("Found total %d violations!\n", numberOfViolations)
				log.Printf("Took %s\n", time.Since(start))
			}

		} else {
			log.Print(err)
		}
	}

	// Print help.
	if *printHelp {
		flag.Usage()
	}
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
