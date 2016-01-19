// The MIT License (MIT)

// Copyright (c) 2015-2016 Christian Bergum Bergersen

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package bblock

import (
	"testing"
	"io/ioutil"
	"go/token"
	"go/parser"
	"go/ast"
	"reflect"
)

type expectedBasicBlock struct {
	Type            BasicBlockType
	Head            ast.Node
	Tail            ast.Node
	SuccessorBlocks []int //Number of block which is successor.
}

//Function which compares generated basic blocks with the corresponding correct basic-block.
//Checks:
// - Number of blocks.
// - Basic-block type.
// - Number of basic-blocks successors.
// - Head and Tail is of correct type.
// - Successor blocks of a basic-block is correct (including order).
func verifyBasicBlocks(expectedBasicBlocks []*BasicBlock, actualBasicBlocks []expectedBasicBlock, t *testing.T) {
	if len(expectedBasicBlocks) != len(actualBasicBlocks) {
		t.Fatalf("Number of basic blocks should be %d, but are %d!\n", len(actualBasicBlocks),
			len(expectedBasicBlocks))
	}

	//Loop through all generated basic-blocks and check if they are similar to the correct once.
	for index, _ := range expectedBasicBlocks {
		if expectedBasicBlocks[index].Type != actualBasicBlocks[index].Type { //Check that basic-block type is correct.
			t.Errorf("Basic block nr. %d should be of type %s, but are of type %s!\n",
				index, actualBasicBlocks[index].Type.String(), expectedBasicBlocks[index].Type.String())
		}
		//Check that length of generate basic-blocks successors are equal correct number of successor blocks.
		if len(expectedBasicBlocks[index].Successor) != len(actualBasicBlocks[index].SuccessorBlocks) {
			t.Fatalf("Number of successors in basic-block nr. %d should be %d, and not %d!\n",
				expectedBasicBlocks[index].Number, len(actualBasicBlocks[index].SuccessorBlocks),
				len(expectedBasicBlocks[index].Successor))
		}
		//Check that Head is of correct type.
		if reflect.TypeOf(expectedBasicBlocks[index].Head) != reflect.TypeOf(actualBasicBlocks[index].Head) {
			t.Errorf("Basic block nr. %d's Head type should be %s, and not %s!\n", expectedBasicBlocks[index].Number,
				reflect.TypeOf(actualBasicBlocks[index].Head), reflect.TypeOf(expectedBasicBlocks[index].Head))
		}
		//Check that Tail is of correct type.
		if reflect.TypeOf(expectedBasicBlocks[index].Tail) != reflect.TypeOf(actualBasicBlocks[index].Tail) {
			t.Errorf("Basic block nr. %d's Tail type should be %s, and not %s!\n", expectedBasicBlocks[index].Number,
				reflect.TypeOf(actualBasicBlocks[index].Tail), reflect.TypeOf(expectedBasicBlocks[index].Tail))
		}

		for i, successorBlock := range expectedBasicBlocks[index].Successor {
			if successorBlock.Number != actualBasicBlocks[index].SuccessorBlocks[i] {
				t.Errorf("Basic block nr. %d's successor block nr. %d should be nr. %d, and not %d!\n",
					index, i, actualBasicBlocks[index].SuccessorBlocks[i], successorBlock.Number)
			}
		}
	}
}

func TestSingleBasicBlock(t *testing.T) {
	srcFile := "./_simple.go"
	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Error finding file: %s\n", srcFile)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		t.Fatalf("Error parsing file: %s\n", srcFile)
	}

	expectedBasicBlocks := GetBasicBlocksFromSourceCode(file)
	actualBasicBlocks := []expectedBasicBlock{
		expectedBasicBlock{Type:PACKAGE_ENTRY, Head:&ast.File{}, Tail:&ast.BasicLit{}, SuccessorBlocks:[]int{1}}, // BB #1
		expectedBasicBlock{Type:FUNCTION_ENTRY, Head:&ast.FuncDecl{}, Tail:&ast.BasicLit{}}, //BB #2
	}
	verifyBasicBlocks(expectedBasicBlocks, actualBasicBlocks, t)
}

func TestIfElseBasicBlock(t *testing.T) {
	srcFile := "./_ifelse.go"
	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Error finding file: %s\n", srcFile)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		t.Fatalf("Error parsing file: %s\n", srcFile)
	}

	expectedBasicBlocks := GetBasicBlocksFromSourceCode(file)
	actualBasicBlocks := []expectedBasicBlock{
		expectedBasicBlock{Type:PACKAGE_ENTRY, Head:&ast.File{}, Tail:&ast.BasicLit{}, SuccessorBlocks:[]int{1}}, // BB #0
		expectedBasicBlock{Type:FUNCTION_ENTRY, Head:&ast.FuncDecl{}, Tail:&ast.BlockStmt{}, SuccessorBlocks:[]int{3, 2}}, //BB #1
		expectedBasicBlock{Type:IF_CONDITION, Head:&ast.IfStmt{}, Tail:&ast.BasicLit{}}, //BB #2
		expectedBasicBlock{Type:ELSE_CONDITION, Head:&ast.BlockStmt{}, Tail:&ast.BasicLit{}}, //BB #3
	}
	verifyBasicBlocks(expectedBasicBlocks, actualBasicBlocks, t)
}

func TestNestedIfElseBasicBlock(t *testing.T) {
	srcFile := "./_nestedifelse.go"
	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Error finding file: %s\n", srcFile)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		t.Fatalf("Error parsing file: %s\n", srcFile)
	}

	expectedBasicBlocks := GetBasicBlocksFromSourceCode(file)
	actualBasicBlocks := []expectedBasicBlock{
		expectedBasicBlock{Type:PACKAGE_ENTRY, Head:&ast.File{}, Tail:&ast.BasicLit{}, SuccessorBlocks:[]int{1}}, // BB #0
		expectedBasicBlock{Type:FUNCTION_ENTRY, Head:&ast.FuncDecl{}, Tail:&ast.BlockStmt{}, SuccessorBlocks:[]int{2, 4}}, //BB #1
		expectedBasicBlock{Type:IF_CONDITION, Head:&ast.IfStmt{}, Tail:&ast.BlockStmt{}, SuccessorBlocks:[]int{3}}, //BB #2
		expectedBasicBlock{Type:IF_CONDITION, Head:&ast.ExprStmt{}, Tail:&ast.BasicLit{}}, //BB #3
		expectedBasicBlock{Type:ELSE_CONDITION, Head:&ast.ExprStmt{}, Tail:&ast.BasicLit{}, SuccessorBlocks:[]int{5, 6}}, //BB #4
		expectedBasicBlock{Type:IF_CONDITION, Head:&ast.ExprStmt{}, Tail:&ast.BasicLit{}}, //BB #5
		expectedBasicBlock{Type:ELSE_CONDITION, Head:&ast.ExprStmt{}, Tail:&ast.BasicLit{}}, //BB #6
	}
	verifyBasicBlocks(expectedBasicBlocks, actualBasicBlocks, t)
}

func TestSwitchBasicBlock(t *testing.T) {
	srcFile := "./_switch.go"
	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Error finding file: %s\n", srcFile)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		t.Fatalf("Error parsing file: %s\n", srcFile)
	}

	expectedBasicBlocks := GetBasicBlocksFromSourceCode(file)
	actualBasicBlocks := []expectedBasicBlock{
		expectedBasicBlock{Type:PACKAGE_ENTRY, Head:&ast.File{}, Tail:&ast.BasicLit{}, SuccessorBlocks:[]int{1}}, // BB #0
		expectedBasicBlock{Type:FUNCTION_ENTRY, Head:&ast.FuncDecl{}, Tail:&ast.BasicLit{}, SuccessorBlocks:[]int{2}}, //BB #1
		expectedBasicBlock{Type:SWITCH_CONDITION, Head:&ast.SwitchStmt{}, Tail:&ast.BlockStmt{}, SuccessorBlocks:[]int{3, 4, 5, 6, 7, 8}}, //BB #2
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #3
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #4
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #5
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #6
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #7
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.Ident{}}, //BB #8
	}
	verifyBasicBlocks(expectedBasicBlocks, actualBasicBlocks, t)
}

func TestNestedSwitchBasicBlock(t *testing.T) {
	srcFile := "./_nestedswitch.go"
	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Error finding file: %s\n", srcFile)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		t.Fatalf("Error parsing file: %s\n", srcFile)
	}

	expectedBasicBlocks := GetBasicBlocksFromSourceCode(file)
	actualBasicBlocks := []expectedBasicBlock{
		expectedBasicBlock{Type:PACKAGE_ENTRY, Head:&ast.File{}, Tail:&ast.BasicLit{}, SuccessorBlocks:[]int{1}}, // BB #0
		expectedBasicBlock{Type:FUNCTION_ENTRY, Head:&ast.FuncDecl{}, Tail:&ast.BasicLit{}, SuccessorBlocks:[]int{2}}, //BB #1
		expectedBasicBlock{Type:SWITCH_CONDITION, Head:&ast.SwitchStmt{}, Tail:&ast.BlockStmt{}, SuccessorBlocks:[]int{3, 4, 9, 10, 11}}, //BB #2
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #3
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #4
		expectedBasicBlock{Type:SWITCH_CONDITION, Head:&ast.SwitchStmt{}, Tail:&ast.BlockStmt{}, SuccessorBlocks:[]int{6, 7, 8}}, //BB #5
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #6
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #7
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.Ident{}}, //BB #8
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #9
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.BasicLit{}}, //BB #10
		expectedBasicBlock{Type:CASE_CLAUSE, Head:&ast.CaseClause{}, Tail:&ast.Ident{}}, //BB #11
	}
	verifyBasicBlocks(expectedBasicBlocks, actualBasicBlocks, t)
}
