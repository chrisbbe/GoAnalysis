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
package bblock_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/chrisbbe/GoAnalysis/analysis/ccomplexity/bblock"
)

// VerifyBasicBlocks checks the list of expected basic-blocks with the list of actual basic-blocks.
func verifyBasicBlocks(expectedBasicBlocks []*bblock.BasicBlock, correctBasicBlocks []*bblock.BasicBlock) error {
	if len(expectedBasicBlocks) != len(correctBasicBlocks) {
		return fmt.Errorf("Number of basic-blocks should be %d, but are %d!\n", len(correctBasicBlocks), len(expectedBasicBlocks))
	}

	//Loop through all generated basic-blocks and check if they are similar to the correct once.
	for index, _ := range expectedBasicBlocks {
		if expectedBasicBlocks[index].Type != correctBasicBlocks[index].Type {
			//Check that basic-block type is correct.
			return fmt.Errorf("Basic block nr. %d should be of type %s, but are of type %s!\n",
				index, correctBasicBlocks[index].Type.String(), expectedBasicBlocks[index].Type.String())
		}

		//Check that length of generate basic-blocks successors are equal correct number of successor blocks.
		if len(expectedBasicBlocks[index].GetSuccessorBlocks()) != len(correctBasicBlocks[index].GetSuccessorBlocks()) {
			return fmt.Errorf("Number of successors in basic-block nr. %d should be %d, and not %d!\n",
				expectedBasicBlocks[index].Number, len(correctBasicBlocks[index].GetSuccessorBlocks()),
				len(expectedBasicBlocks[index].GetSuccessorBlocks()))
		}

		//Check that basic block starts at right line.
		if expectedBasicBlocks[index].EndLine != correctBasicBlocks[index].EndLine {
			return fmt.Errorf("Basic block nr. %d should end at line number %d, and not %d!\n", expectedBasicBlocks[index].Number,
				correctBasicBlocks[index].EndLine, expectedBasicBlocks[index].EndLine)
		}

		//Check that that basic-block has correct successor blocks, and their order.
		for i, successorBlock := range expectedBasicBlocks[index].GetSuccessorBlocks() {
			if successorBlock.Number != correctBasicBlocks[index].GetSuccessorBlocks()[i].Number {
				return fmt.Errorf("Basic block nr. %d's successor block nr. %d should be nr. %d, and not %d!\n",
					index, i, correctBasicBlocks[index].GetSuccessorBlocks()[i].Number, successorBlock.Number)
			}
		}

	}
	return nil
}

func TestEmptyFunctionBasicBlock(t *testing.T) {
	sourceFile, err := ioutil.ReadFile("./testcode/_emptyFunction.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(sourceFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 29)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 31)

	BB0.AddSuccessorBlock(BB1)

	correctBasicBlocks := []*bblock.BasicBlock{BB0, BB1}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestSingleBasicBlock(t *testing.T) {
	sourceFile, err := ioutil.ReadFile("./testcode/_simple.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(sourceFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 31)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 33)

	BB0.AddSuccessorBlock(BB1)

	correctBasicBlocks := []*bblock.BasicBlock{BB0, BB1}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestIfElseBasicBlock(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_ifelse.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 29)
	BB1 := bblock.NewBasicBlock(1, bblock.IF_CONDITION, 33)
	BB2 := bblock.NewBasicBlock(2, bblock.ELSE_CONDITION, 35)
	BB3 := bblock.NewBasicBlock(3, bblock.ELSE_BODY, 38)
	BB4 := bblock.NewBasicBlock(4, bblock.RETURN_STMT, 39)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB3)
	BB2.AddSuccessorBlock(BB4)
	BB3.AddSuccessorBlock(BB4)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestNestedIfElseBasicBlock(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_nestedifelse.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 5)
	BB1 := bblock.NewBasicBlock(1, bblock.IF_CONDITION, 7)
	BB2 := bblock.NewBasicBlock(2, bblock.IF_CONDITION, 10)
	BB3 := bblock.NewBasicBlock(3, bblock.ELSE_CONDITION, 12)
	BB4 := bblock.NewBasicBlock(4, bblock.ELSE_BODY, 14)
	BB5 := bblock.NewBasicBlock(5, bblock.ELSE_CONDITION, 15)
	BB6 := bblock.NewBasicBlock(6, bblock.ELSE_BODY, 17)
	BB7 := bblock.NewBasicBlock(7, bblock.RETURN_STMT, 20)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB6)
	BB2.AddSuccessorBlock(BB3, BB4)
	BB3.AddSuccessorBlock(BB7)
	BB4.AddSuccessorBlock(BB7)
	BB5.AddSuccessorBlock(BB7)
	BB6.AddSuccessorBlock(BB7)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4, BB5, BB6, BB7,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestLooperBasicBlock(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_looper.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 3)
	BB1 := bblock.NewBasicBlock(1, bblock.FOR_STATEMENT, 5)
	BB2 := bblock.NewBasicBlock(2, bblock.FOR_BODY, 7)
	BB3 := bblock.NewBasicBlock(3, bblock.RETURN_STMT, 9)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB3)
	BB2.AddSuccessorBlock(BB1)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleSwitchBasicBlock(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_simpleSwitch.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 31)
	BB1 := bblock.NewBasicBlock(1, bblock.SWITCH_STATEMENT, 35)
	BB2 := bblock.NewBasicBlock(2, bblock.CASE_CLAUSE, 38)
	BB3 := bblock.NewBasicBlock(3, bblock.CASE_CLAUSE, 40)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 42)
	BB5 := bblock.NewBasicBlock(5, bblock.CASE_CLAUSE, 44)
	BB6 := bblock.NewBasicBlock(6, bblock.RETURN_STMT, 46)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB3, BB4, BB5, BB6)
	BB2.AddSuccessorBlock(BB6)
	BB3.AddSuccessorBlock(BB6)
	BB4.AddSuccessorBlock(BB6)
	BB5.AddSuccessorBlock(BB6)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4, BB5, BB6,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestSwitchBasicBlock(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_switch.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 31)
	BB1 := bblock.NewBasicBlock(1, bblock.SWITCH_STATEMENT, 34)
	BB2 := bblock.NewBasicBlock(2, bblock.CASE_CLAUSE, 37)
	BB3 := bblock.NewBasicBlock(3, bblock.CASE_CLAUSE, 40)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 42)
	BB5 := bblock.NewBasicBlock(5, bblock.CASE_CLAUSE, 44)
	BB6 := bblock.NewBasicBlock(6, bblock.RETURN_STMT, 47)
	BB7 := bblock.NewBasicBlock(7, bblock.CASE_CLAUSE, 49)
	BB8 := bblock.NewBasicBlock(8, bblock.RETURN_STMT, 51)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB3, BB4, BB5, BB6, BB7, BB8)
	BB2.AddSuccessorBlock(BB8)
	BB3.AddSuccessorBlock(BB8)
	BB4.AddSuccessorBlock(BB8)
	BB5.AddSuccessorBlock(BB8)
	BB7.AddSuccessorBlock(BB8)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4, BB5, BB6, BB7, BB8,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestReturnSwitcherBasicBlock(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_returnSwitcher.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 31)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 34)
	BB2 := bblock.NewBasicBlock(2, bblock.FUNCTION_ENTRY, 36)
	BB3 := bblock.NewBasicBlock(3, bblock.SWITCH_STATEMENT, 37)
	BB4 := bblock.NewBasicBlock(4, bblock.RETURN_STMT, 39)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 41)
	BB6 := bblock.NewBasicBlock(6, bblock.RETURN_STMT, 43)
	BB7 := bblock.NewBasicBlock(7, bblock.RETURN_STMT, 45)
	BB8 := bblock.NewBasicBlock(8, bblock.RETURN_STMT, 47)
	BB9 := bblock.NewBasicBlock(9, bblock.RETURN_STMT, 49)
	BB10 := bblock.NewBasicBlock(10, bblock.RETURN_STMT, 51)
	BB11 := bblock.NewBasicBlock(11, bblock.RETURN_STMT, 53)
	BB12 := bblock.NewBasicBlock(12, bblock.RETURN_STMT, 55)
	BB13 := bblock.NewBasicBlock(13, bblock.RETURN_STMT, 57)
	BB14 := bblock.NewBasicBlock(14, bblock.RETURN_STMT, 59)
	BB15 := bblock.NewBasicBlock(15, bblock.RETURN_STMT, 61)
	BB16 := bblock.NewBasicBlock(16, bblock.RETURN_STMT, 63)

	// Function main.
	BB0.AddSuccessorBlock(BB1)

	// Function monthNumberToString.
	BB2.AddSuccessorBlock(BB3)
	BB3.AddSuccessorBlock(BB4, BB5, BB6, BB7, BB8, BB9, BB10, BB11, BB12, BB13, BB14, BB15, BB16)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4, BB5, BB6, BB7, BB8, BB9, BB10, BB11, BB12, BB13, BB14, BB15, BB16,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestNestedSwitchBasicBlock(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_nestedswitch.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 33)
	BB1 := bblock.NewBasicBlock(1, bblock.SWITCH_STATEMENT, 37)
	BB2 := bblock.NewBasicBlock(2, bblock.CASE_CLAUSE, 40)
	BB3 := bblock.NewBasicBlock(3, bblock.SWITCH_STATEMENT, 43)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 46)
	BB5 := bblock.NewBasicBlock(5, bblock.CASE_CLAUSE, 48)
	BB6 := bblock.NewBasicBlock(6, bblock.CASE_CLAUSE, 50)
	BB7 := bblock.NewBasicBlock(7, bblock.CASE_CLAUSE, 53)
	BB8 := bblock.NewBasicBlock(8, bblock.CASE_CLAUSE, 55)
	BB9 := bblock.NewBasicBlock(9, bblock.CASE_CLAUSE, 57)
	BB10 := bblock.NewBasicBlock(10, bblock.RETURN_STMT, 59)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4, BB5, BB6, BB7, BB8, BB9, BB10,
	}

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB3, BB7, BB8, BB9, BB10)
	BB2.AddSuccessorBlock(BB10)
	BB3.AddSuccessorBlock(BB4, BB5, BB6, BB10)
	BB4.AddSuccessorBlock(BB10)
	BB5.AddSuccessorBlock(BB10)
	BB6.AddSuccessorBlock(BB10)
	BB7.AddSuccessorBlock(BB10)
	BB8.AddSuccessorBlock(BB10)
	BB9.AddSuccessorBlock(BB10)

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestTypeSwitchBasicBlock(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_typeswitch.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Error(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 31)
	BB1 := bblock.NewBasicBlock(1, bblock.SWITCH_STATEMENT, 36)
	BB2 := bblock.NewBasicBlock(2, bblock.CASE_CLAUSE, 39)
	BB3 := bblock.NewBasicBlock(3, bblock.CASE_CLAUSE, 41)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 43)
	BB5 := bblock.NewBasicBlock(5, bblock.CASE_CLAUSE, 45)
	BB6 := bblock.NewBasicBlock(6, bblock.CASE_CLAUSE, 47)
	BB7 := bblock.NewBasicBlock(7, bblock.RETURN_STMT, 50)

	BB0.AddSuccessorBlock(BB1)

	BB1.AddSuccessorBlock(BB2)
	BB1.AddSuccessorBlock(BB3)
	BB1.AddSuccessorBlock(BB4)
	BB1.AddSuccessorBlock(BB5)
	BB1.AddSuccessorBlock(BB6)

	BB2.AddSuccessorBlock(BB7)
	BB3.AddSuccessorBlock(BB7)
	BB4.AddSuccessorBlock(BB7)
	BB5.AddSuccessorBlock(BB7)
	BB6.AddSuccessorBlock(BB7)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4, BB5, BB6, BB7,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleLooperSwitch(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_simplelooperswitch.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Error(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 5)
	BB1 := bblock.NewBasicBlock(1, bblock.FOR_STATEMENT, 7)
	BB2 := bblock.NewBasicBlock(2, bblock.SWITCH_STATEMENT, 8)
	BB3 := bblock.NewBasicBlock(3, bblock.CASE_CLAUSE, 10)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 12)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 15)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB5)
	BB2.AddSuccessorBlock(BB1, BB3, BB4)
	BB3.AddSuccessorBlock(BB1)
	BB4.AddSuccessorBlock(BB1)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4, BB5,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestSelectBasicBlock(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_select.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 35)
	BB1 := bblock.NewBasicBlock(1, bblock.GO_STATEMENT, 38)
	BB2 := bblock.NewBasicBlock(2, bblock.FOR_STATEMENT, 45)
	BB3 := bblock.NewBasicBlock(3, bblock.SELECT_STATEMENT, 46)
	BB4 := bblock.NewBasicBlock(4, bblock.COMM_CLAUSE, 48)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 51)
	BB6 := bblock.NewBasicBlock(6, bblock.RETURN_STMT, 54)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2)
	BB2.AddSuccessorBlock(BB3, BB6)
	BB3.AddSuccessorBlock(BB2, BB4, BB5)
	BB4.AddSuccessorBlock(BB2)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4, BB5, BB6,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}

func TestGreatestCommonDivisor(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_gcd.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 31)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 34)
	BB2 := bblock.NewBasicBlock(2, bblock.FUNCTION_ENTRY, 36)
	BB3 := bblock.NewBasicBlock(3, bblock.FOR_STATEMENT, 37)
	BB4 := bblock.NewBasicBlock(4, bblock.FOR_BODY, 39)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 40)

	BB0.AddSuccessorBlock(BB1)
	BB2.AddSuccessorBlock(BB3)
	BB3.AddSuccessorBlock(BB4)
	BB3.AddSuccessorBlock(BB5)
	BB4.AddSuccessorBlock(BB3)

	correctBasicBlocks := []*bblock.BasicBlock{
		BB0, BB1, BB2, BB3, BB4, BB5,
	}

	if err := verifyBasicBlocks(expectedBasicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
}
