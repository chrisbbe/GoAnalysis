// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package bblock_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/chrisbbe/GoAnalysis/analyzer/ccomplexity/bblock"
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
	sourceFile, err := ioutil.ReadFile("./testcode/_emptyfunction.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(sourceFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 6)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 8)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 11)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 6)
	BB1 := bblock.NewBasicBlock(1, bblock.IF_CONDITION, 10)
	BB2 := bblock.NewBasicBlock(2, bblock.ELSE_CONDITION, 12)
	BB3 := bblock.NewBasicBlock(3, bblock.ELSE_BODY, 15)
	BB4 := bblock.NewBasicBlock(4, bblock.RETURN_STMT, 16)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.IF_CONDITION, 11)
	BB2 := bblock.NewBasicBlock(2, bblock.IF_CONDITION, 15)
	BB3 := bblock.NewBasicBlock(3, bblock.ELSE_CONDITION, 18)
	BB4 := bblock.NewBasicBlock(4, bblock.ELSE_BODY, 21)
	BB5 := bblock.NewBasicBlock(5, bblock.ELSE_CONDITION, 22)
	BB6 := bblock.NewBasicBlock(6, bblock.ELSE_BODY, 25)
	BB7 := bblock.NewBasicBlock(7, bblock.RETURN_STMT, 28)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.FOR_STATEMENT, 11)
	BB2 := bblock.NewBasicBlock(2, bblock.FOR_BODY, 14)
	BB3 := bblock.NewBasicBlock(3, bblock.RETURN_STMT, 16)

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
	srcFile, err := ioutil.ReadFile("./testcode/_simpleswitch.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.SWITCH_STATEMENT, 12)
	BB2 := bblock.NewBasicBlock(2, bblock.CASE_CLAUSE, 15)
	BB3 := bblock.NewBasicBlock(3, bblock.CASE_CLAUSE, 17)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 19)
	BB5 := bblock.NewBasicBlock(5, bblock.CASE_CLAUSE, 21)
	BB6 := bblock.NewBasicBlock(6, bblock.RETURN_STMT, 23)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.SWITCH_STATEMENT, 12)
	BB2 := bblock.NewBasicBlock(2, bblock.CASE_CLAUSE, 15)
	BB3 := bblock.NewBasicBlock(3, bblock.CASE_CLAUSE, 18)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 20)
	BB5 := bblock.NewBasicBlock(5, bblock.CASE_CLAUSE, 22)
	BB6 := bblock.NewBasicBlock(6, bblock.RETURN_STMT, 25)
	BB7 := bblock.NewBasicBlock(7, bblock.CASE_CLAUSE, 27)
	BB8 := bblock.NewBasicBlock(8, bblock.RETURN_STMT, 29)

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
	srcFile, err := ioutil.ReadFile("./testcode/_returnswitcher.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedBasicBlocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 11)
	BB2 := bblock.NewBasicBlock(2, bblock.FUNCTION_ENTRY, 13)
	BB3 := bblock.NewBasicBlock(3, bblock.SWITCH_STATEMENT, 14)
	BB4 := bblock.NewBasicBlock(4, bblock.RETURN_STMT, 16)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 18)
	BB6 := bblock.NewBasicBlock(6, bblock.RETURN_STMT, 20)
	BB7 := bblock.NewBasicBlock(7, bblock.RETURN_STMT, 22)
	BB8 := bblock.NewBasicBlock(8, bblock.RETURN_STMT, 24)
	BB9 := bblock.NewBasicBlock(9, bblock.RETURN_STMT, 26)
	BB10 := bblock.NewBasicBlock(10, bblock.RETURN_STMT, 28)
	BB11 := bblock.NewBasicBlock(11, bblock.RETURN_STMT, 30)
	BB12 := bblock.NewBasicBlock(12, bblock.RETURN_STMT, 32)
	BB13 := bblock.NewBasicBlock(13, bblock.RETURN_STMT, 34)
	BB14 := bblock.NewBasicBlock(14, bblock.RETURN_STMT, 36)
	BB15 := bblock.NewBasicBlock(15, bblock.RETURN_STMT, 38)
	BB16 := bblock.NewBasicBlock(16, bblock.RETURN_STMT, 40)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 10)
	BB1 := bblock.NewBasicBlock(1, bblock.SWITCH_STATEMENT, 14)
	BB2 := bblock.NewBasicBlock(2, bblock.CASE_CLAUSE, 17)
	BB3 := bblock.NewBasicBlock(3, bblock.SWITCH_STATEMENT, 20)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 23)
	BB5 := bblock.NewBasicBlock(5, bblock.CASE_CLAUSE, 25)
	BB6 := bblock.NewBasicBlock(6, bblock.CASE_CLAUSE, 27)
	BB7 := bblock.NewBasicBlock(7, bblock.CASE_CLAUSE, 30)
	BB8 := bblock.NewBasicBlock(8, bblock.CASE_CLAUSE, 32)
	BB9 := bblock.NewBasicBlock(9, bblock.CASE_CLAUSE, 34)
	BB10 := bblock.NewBasicBlock(10, bblock.RETURN_STMT, 36)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.SWITCH_STATEMENT, 14)
	BB2 := bblock.NewBasicBlock(2, bblock.CASE_CLAUSE, 17)
	BB3 := bblock.NewBasicBlock(3, bblock.CASE_CLAUSE, 19)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 21)
	BB5 := bblock.NewBasicBlock(5, bblock.CASE_CLAUSE, 23)
	BB6 := bblock.NewBasicBlock(6, bblock.CASE_CLAUSE, 25)
	BB7 := bblock.NewBasicBlock(7, bblock.RETURN_STMT, 28)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.FOR_STATEMENT, 11)
	BB2 := bblock.NewBasicBlock(2, bblock.SWITCH_STATEMENT, 13)
	BB3 := bblock.NewBasicBlock(3, bblock.CASE_CLAUSE, 15)
	BB4 := bblock.NewBasicBlock(4, bblock.CASE_CLAUSE, 17)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 20)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 12)
	BB1 := bblock.NewBasicBlock(1, bblock.GO_STATEMENT, 16)
	BB2 := bblock.NewBasicBlock(2, bblock.FOR_STATEMENT, 24)
	BB3 := bblock.NewBasicBlock(3, bblock.SELECT_STATEMENT, 26)
	BB4 := bblock.NewBasicBlock(4, bblock.COMM_CLAUSE, 28)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 31)
	BB6 := bblock.NewBasicBlock(6, bblock.RETURN_STMT, 34)

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

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 12)
	BB2 := bblock.NewBasicBlock(2, bblock.FUNCTION_ENTRY, 14)
	BB3 := bblock.NewBasicBlock(3, bblock.FOR_STATEMENT, 16)
	BB4 := bblock.NewBasicBlock(4, bblock.FOR_BODY, 19)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 20)

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
