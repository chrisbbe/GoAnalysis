package bblock

import (
	"testing"
	"io/ioutil"
	"go/token"
	"go/parser"
)

func TestIfElseBasicBlock(t *testing.T) {

	srcFile := "./_ifelse.go"

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Error finding file %s!\n", srcFile)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	basicBlocks := GetBasicBlocksFromSourceCode(fset, file)

	numberOfBlocks := 5

	if len(basicBlocks) != numberOfBlocks {
		t.Errorf("Number of basic blocks should be %d, but are %d!\n", numberOfBlocks, len(basicBlocks))
	}

}

func TestAnother(t *testing.T) {
	srcFile := "./_another.go"

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Error finding file %s!\n", srcFile)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	basicBlocks := GetBasicBlocksFromSourceCode(fset, file)
	numberOfBlocks := 9

	if len(basicBlocks) != numberOfBlocks {
		t.Errorf("Number of basic blocks should be %d, but are %d!\n", numberOfBlocks, len(basicBlocks))
	}
}
