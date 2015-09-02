package main

import (
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"go/parser"
	"go/token"
	"go/ast"
	"io"
	"./stack"
)

func main() {
	srcFile, err := filenamesFromCommandLine()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("srcFile = %s\n", srcFile)

	src, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error:\n")
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	filename := "ast.gv"

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}

	writeLine("digraph AST {\n", f)

	fv := new(Visitor)
	fv.outputFile = f
	ast.Walk(fv, file)

	writeLine("}\n", f)
	f.Close()
	// Print the AST.
	//ast.Print(fset, file)
}

func filenamesFromCommandLine() (srcFilename string, err error) {
	if len(os.Args) > 2 && os.Args[1] == "-s" {
		return os.Args[2], nil
	}
	err = fmt.Errorf("usage: %s -s goSource.go\n", filepath.Base(os.Args[0]))
	return "", err
}

func writeLine(line string, f io.Writer) {
	n, err := io.WriteString(f, line)
	if(err != nil) {
		fmt.Println(n, err)
	}
}

type Visitor struct {
	nodeStack stack.Stack
	outputFile io.Writer
}

func (v *Visitor) Visit(node ast.Node) (w ast.Visitor)  {

	if(node != nil) {
		//fmt.Printf("Node: %T\n", node)

		val, _ := v.nodeStack.Top()
		//fmt.Printf("Father is %T\n\n", val)

		if(val != nil) {
			var l = ""

			switch t := node.(type) {

			case *ast.BasicLit:
				l = fmt.Sprintf("\"%T\" -> \"%T: %s\";\n", val, node, t.Kind)

			default:
				l = fmt.Sprintf("\"%T\" -> \"%T\";\n", val, node)
			}

			writeLine(l, v.outputFile)
		}
	}

	if(node != nil) {
		v.nodeStack.Push(node)
	} else {
		v.nodeStack.Pop()
	}

	return v
}
