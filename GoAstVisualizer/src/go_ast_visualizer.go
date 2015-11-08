// This code is part of a master thesis written by
// Christian Bergum Bergersen at University of Oslo,
// Faculty of Mathematics and Natural Sciences,
// Departments of Computer Science.

// TODO: Add some text about license for the code.

// This program utilizes the standard ast package in Go
// and the corresponding Visit() method to inspect the
// abstract syntax three for Go programs specified as
// input, the program then writes dotty codes for each
// node out to a separate file which can be compiled
// for graphical visualization of the abstract syntax tree.

package main

import (
	"./stack"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	srcFile, err := getFilenameFromCommandLine()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("## Go Abstract Syntax Tree Visualizer ##")
	fmt.Printf("Go version: %s\n", runtime.Version())

	src, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error:\n")
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	inFilename := strings.Split(filepath.Base(srcFile), ".go")
	outFile := strings.Join(strings.Split(filepath.Base(srcFile), ".go"), ".gv")
	fmt.Printf("Outfile: %s\n", outFile)

	dottyFile, err := os.Create(outFile)
	if err != nil {
		fmt.Println(err)
	}

	writeLineToFile("digraph AST {\n", dottyFile)

	fv := visitor{fileSet:fset, outputFile:dottyFile}
	ast.Walk(&fv, file)

	writeLineToFile("}\n", dottyFile)
	dottyFile.Close()

	fmt.Printf("Run: $ dot -Tpdf %s -o %s.pdf\nto create PDF of abstract syntax tree.\n", outFile, inFilename[0])
}

type visitor struct {
	fileSet    *token.FileSet
	nodeStack  stack.Stack
	outputFile io.Writer
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		val, _ := v.nodeStack.Top()

		if val != nil {
			var line = ""
			tt := val.(ast.Node)
			switch t := node.(type) {

			case *ast.Ident:
				line = fmt.Sprintf("\t\"%T (Line: %d)\" -> \"%T: %s (Line: %d)\";\n", val,
					getLineNumberInSourceCode(v.fileSet, tt.Pos()), node, t.Name,
					getLineNumberInSourceCode(v.fileSet, t.Pos()))

			case *ast.BasicLit:
				line = fmt.Sprintf("\t\"%T (Line: %d)\" -> \"%T: %s (Line: %d)\";\n", val,
					getLineNumberInSourceCode(v.fileSet, tt.Pos()), node, strExtract(t.Value, "\""),
					getLineNumberInSourceCode(v.fileSet, t.Pos()))

			default:
				line = fmt.Sprintf("\t\"%T (Line: %d)\" -> \"%T (Line: %d)\";\n", val,
					getLineNumberInSourceCode(v.fileSet, tt.Pos()), node,
					getLineNumberInSourceCode(v.fileSet, node.Pos()))
			}
			writeLineToFile(line, v.outputFile)
		}
		v.nodeStack.Push(node) //Push node on stack, we will go further down.
	} else {
		v.nodeStack.Pop() //Pop node from stack, going one level up to parent.
	}

	return v
}

func strExtract(line string, delimiter string) string {
	splittedString := strings.Split(line, delimiter)
	if len(splittedString) == 1 {
		return ""
	}
	splittedString = strings.Split(splittedString[1], delimiter)
	return splittedString[0]
}

func getFilenameFromCommandLine() (srcFilename string, err error) {
	if len(os.Args) > 2 && os.Args[1] == "-s" {
		return os.Args[2], nil
	}
	err = fmt.Errorf("Usage: %s -s go_source.go\n", filepath.Base(os.Args[0]))
	return "", err
}

func writeLineToFile(line string, f io.Writer) {
	n, err := io.WriteString(f, line)
	if err != nil {
		fmt.Println(n, err)
	}
}

func getLineNumberInSourceCode(fileSet *token.FileSet, position token.Pos) (line int) {
	tokenFile := fileSet.File(position)
	return tokenFile.Line(position)
}
