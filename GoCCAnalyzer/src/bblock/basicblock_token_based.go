package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
)

func main() {
	src, err := ioutil.ReadFile("GoCCAnalyzer/src/bblock/test.go")
	if err != nil {
		fmt.Println("Error: Cloudnt find file")
	}

	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	counter := 0
	inElse := false

	blockCounter := 0

	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok == token.RETURN || tok == token.IF || tok == token.ELSE {
			fmt.Printf("\n*********** BLOCK NO. %d *********************\n", blockCounter)
			blockCounter++
			if tok == token.ELSE {
				inElse = true
			}
		}

		fmt.Printf("%s\t%q\n", tok, lit)

		if inElse {
			if tok == token.LBRACE {
				counter++
			}
			if tok == token.RBRACE {
				counter--
				if counter == 0 {
					fmt.Printf("\n*********** BLOCK NO. %d *********************\n", blockCounter)
					blockCounter++
					inElse = false
				}
			}
		}
	}
}
