// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"encoding/json"
	"flag"
	"github.com/chrisbbe/GoAnalysis/analyzer/globalvars"
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

	// Option dir selected.
	if len(*sourceRootDir) > 0 {
		start := time.Now()

		goPackageViolations, err := linter.DetectViolations(*sourceRootDir) // Start the analysis.
		if err != nil {
			log.Fatal(err)
		}
		timeUsed := time.Since(start)

		// Direct output to console as JSON.
		if *jsonOutput && len(goPackageViolations) > 0 {

			// Aggregate GoFiles to JSON marshalling.
			violationsMap := make(map[string]*linter.GoFile)
			for _, goPackage := range goPackageViolations {
				for _, goFile := range goPackage.Violations {
					if gf, ok := violationsMap[goFile.FilePath]; ok {
						gf.Violations = append(gf.Violations, goFile.Violations...)
					} else {
						violationsMap[goFile.FilePath] = goFile
					}
				}
			}

			// Convert to list for JSON marshalling.
			var violations []*linter.GoFile
			for _, value := range violationsMap {
				violations = append(violations, value)
			}

			json, err := json.MarshalIndent(violations, "", "\t")
			if err != nil {
				log.Fatal(err)
			}
			if _, err := os.Stdout.Write(json); err != nil {
				panic(err)
			}

		} else {
			log.SetOutput(os.Stdout) // We want to send output to stdout, instead of Stderr.
			numberOfViolations := 0
			linesOfCode := 0
			linesOfComments := 0

			// Nicely print the output to the console.
			log.Println("-----------------------------------------------------------------------------------------------")
			for _, goPackage := range goPackageViolations {
				log.Printf("PACKAGE: %s (%s)", goPackage.Pack.Name, goPackage.Path)

				for _, goFile := range goPackage.Violations {
					log.Printf("\tViolations in %s :", filepath.Base(goFile.FilePath))
					for i, vio := range goFile.Violations {
						log.Printf("\t\t%d) %s\n", i, vio)
					}
					linesOfCode += goFile.LinesOfCode
					linesOfComments += goFile.LinesOfComments
				}

				numberOfViolations += len(goPackage.Violations)
				log.Println("-----------------------------------------------------------------------------------------------")
			}
			log.Println("## ANALYSIS SUMMARY ##")
			log.Printf("Total %d violations found!\n", numberOfViolations)
			log.Printf("Total number of Go files: %d\n", countGoFiles(*sourceRootDir))
			log.Printf("Total lines of code (LOC): %d\n", linesOfCode)
			log.Printf("Total lines of comments: %d\n", linesOfComments)
			log.Printf("Total time used: %s\n", timeUsed)
			log.Printf("For rule details: %s\n", globalvars.WIKI_PAGE)
		}

	}

	// Print help.
	if *printHelp {
		flag.Usage()
	}
}

// getGoFiles searches recursively for .go files in the searchDir path, returning the absolute path to the files.
func countGoFiles(searchDir string) (counter int) {
	filepath.Walk(searchDir, func(pat string, file os.FileInfo, err error) error {

		if file != nil && !file.IsDir() && file.Mode().IsRegular() && strings.EqualFold(path.Ext(file.Name()), ".go") {
			// Only add if regular Go source-code file.
			counter++
		}
		return nil
	})
	return
}
