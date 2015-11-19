package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

type basicBlock struct {
	number   int
	fromLine int
	toLine   int
}

type visitor struct {
	fileSet           *token.FileSet
	blockCounter      int
	basicBlocks       []*basicBlock
	leaderSourceLine  int
	currentSourceLine int
}

func main() {
	srcFile := "/uio/hume/student-u85/chrisbbe/IdeaProjects/GoThesis/GoCCAnalyzer/src/bblock/test_src.go"

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error finding file\n")
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	//ast.Print(fset, file)
	fv := visitor{fileSet: fset, blockCounter: 0, leaderSourceLine: -1}
	ast.Walk(&fv, file)
	fmt.Println("### PRINTING RESULT ###")

	for _, bb := range fv.basicBlocks {
		fmt.Printf("################## BLOCK NR. %d (%d - %d) ##################\n", bb.number, bb.fromLine, bb.toLine)
	}
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch t := node.(type) {

		case *ast.IfStmt:
			fmt.Printf("################## BLOCK NR. %d ##################\n", v.blockCounter)
			line := getSourceCodeLine(v.fileSet, t.Pos())
			fmt.Printf("%T (%d)\n", t, line)

			bb := basicBlock{number: v.blockCounter, fromLine: v.leaderSourceLine, toLine: getSourceCodeLine(v.fileSet, t.Pos())}
			v.basicBlocks = append(v.basicBlocks, &bb)
			v.leaderSourceLine = -1
			v.blockCounter++

		case *ast.SwitchStmt:
			fmt.Printf("################## BLOCK NR. %d ##################\n", v.blockCounter)
			line := getSourceCodeLine(v.fileSet, t.Pos())
			fmt.Printf("%T (%d)\n", t, line)

			bb := basicBlock{number: v.blockCounter, fromLine: v.leaderSourceLine, toLine: getSourceCodeLine(v.fileSet, t.Pos())}
			v.basicBlocks = append(v.basicBlocks, &bb)
			v.leaderSourceLine = -1
			v.blockCounter++

		case *ast.CaseClause:
			fmt.Printf("################## BLOCK NR. %d (%d) ##################\n", v.blockCounter, getSourceCodeLine(v.fileSet, t.Pos()))
			line := getSourceCodeLine(v.fileSet, t.Pos())
			fmt.Printf("%T (%d)\n", t, line)

			bb := basicBlock{number: v.blockCounter, fromLine: v.leaderSourceLine, toLine: getSourceCodeLine(v.fileSet, t.Pos())}
			v.basicBlocks = append(v.basicBlocks, &bb)
			v.leaderSourceLine = -1
			v.blockCounter++

		default:
			if v.leaderSourceLine == -1 {
				v.leaderSourceLine = getSourceCodeLine(v.fileSet, t.Pos())
			}
			line := getSourceCodeLine(v.fileSet, t.Pos())
			//Create start BBlock
			if len(v.basicBlocks) == 0 {
				fmt.Printf("################## BLOCK NR. %d (%d) ##################\n", v.blockCounter, getSourceCodeLine(v.fileSet, t.Pos()))
				bb := basicBlock{number: v.blockCounter, fromLine: v.leaderSourceLine, toLine: getSourceCodeLine(v.fileSet, t.Pos())}
				v.basicBlocks = append(v.basicBlocks, &bb)
				v.leaderSourceLine = -1
				v.blockCounter++
			}
			fmt.Printf("%T (%d)\n", t, line)

		}
	}
	return v
}

func getSourceCodeLine(fileSet *token.FileSet, pos token.Pos) int {
	tokenFile := fileSet.File(pos)
	return tokenFile.Line(pos)
}
