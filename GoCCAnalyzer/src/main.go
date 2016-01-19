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
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"go/token"
	"go/parser"
	"github.com/chrisbbe/GoAnalysis/GoCCAnalyzer/src/graph"
	"github.com/chrisbbe/GoAnalysis/GoCCAnalyzer/src/bblock"
)

func main() {
	generateControlFlowGraph()
	//getBasicBlocks()
	//generateGraph()
}

func generateGraph() {
	srcFile := "../directedgraph.txt"

	file, err := os.Open(srcFile)
	if err != nil {
		fmt.Printf("Error opening file %s!\n", srcFile)
		os.Exit(1)
	}
	defer file.Close()

	g := graph.New()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " -> ")
		left := graph.Node{Value: line[0]}
		right := graph.Node{Value: line[1]}
		g.InsertNode(&left, &right)
	}

	fmt.Println("### DFS Printout ###")
	for _, node := range g.GetDFS() {
		fmt.Printf("%s\n", node.Value)
	}
}

func getBasicBlocks() {
	srcFile := "./bblock/_nestedifelse.go"

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("File does not exist: %s!\n", srcFile)
		os.Exit(1)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	basicBlocks := bblock.GetBasicBlocksFromSourceCode(file)

	for _, bb := range basicBlocks {
		fmt.Printf("******** BLOCK NR. %d (%s) ******************\n", bb.Number, bb.Type.String())

		fmt.Printf("\t\tHead: %T, Tail: %T\n", bb.Head, bb.Tail)

		for _, s := range bb.Successor {
			fmt.Printf("\t\t- BLOCK NR. %d (%s)\n", s.Number, s.Type.String())
		}
	}
}

type controlFlowNode struct {
	Value string
	bb    *bblock.BasicBlock
}

func generateControlFlowGraph() {
	srcFile := "../codeexamples/_ifelse.go"

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("File does not exist: %s!\n", srcFile)
		os.Exit(1)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	controlFlowGraph := graph.New()

	for _, basicBlock := range bblock.GetBasicBlocksFromSourceCode(file) {
		if len(basicBlock.Successor) != 0 {
			value := fmt.Sprintf("%s (%d)", basicBlock.Type.String(), basicBlock.Number)
			leftNode := graph.Node{Value:controlFlowNode{Value:value, bb:basicBlock}}

			for _, successorBlock := range basicBlock.Successor {
				value := fmt.Sprintf("%s (%d)", successorBlock.Type.String(), successorBlock.Number)
				rightNode := graph.Node{Value:controlFlowNode{Value:value, bb:successorBlock}}
				controlFlowGraph.InsertNode(&leftNode, &rightNode)
			}
		}
	}

	fmt.Printf("Number of nodes in graph: %d\n", controlFlowGraph.GetNumberOfNodes())
	fmt.Printf("Number of edges in graph: %d\n", controlFlowGraph.GetNumberOfEdges())
	cyclomaticComplexity := controlFlowGraph.GetNumberOfEdges() - controlFlowGraph.GetNumberOfNodes() + 2
	fmt.Printf("Cyclomatic complexity: %d\n", cyclomaticComplexity)

	fmt.Println("\n* Depth First Search *")
	for _, node := range controlFlowGraph.GetDFS() {
		fmt.Printf("%s\n", node.Value.(controlFlowNode).Value)
	}

}
