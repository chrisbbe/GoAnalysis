//This package provides the ability to parse Go source code
//and detect common Go mistakes like:
//	- Usage of loop-iterator variables in goroutines.
//
//Author: Christian Bergum Bergersen (chrisbbe@ifi.uio.no)
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/chrisbbe/GoThesis/GoCommonMistakeDetector/src/mistake"
)

func main() {
	srcFile, err := getFilenameFromCommandLine()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("This is the Go commom mistake detector... Looking for mistakes...\n")
	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error:\n")
	}

	result, err := mistake.FindCommonMistakes(sourceFile)
	if err != nil {
		fmt.Println(err)
	}

	for _, value := range result {
		if value.MistakeType == mistake.RACE_CONDITION {
			fmt.Printf("Warning: Potential race-condition found on line %d.\n", value.LineInSourceFile)
		}
	}
}

func getFilenameFromCommandLine() (srcFilename string, err error) {
	if len(os.Args) > 2 && os.Args[1] == "-s" {
		return os.Args[2], nil
	}
	err = fmt.Errorf("Usage: %s -s <go_source.go>\n", filepath.Base(os.Args[0]))
	return "", err
}
