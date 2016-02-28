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

// Package 'ccomplexity' provides functions to measure cyclomatic complexity
// of Go source code files, both cyclomatic complexity per function and file is
// supported.
package ccomplexity

import (
	"github.com/chrisbbe/GoAnalysis/GoCCAnalyzer/src/graph"
	"github.com/chrisbbe/GoAnalysis/GoCCAnalyzer/src/bblock"
	"go/token"
	"go/parser"
	"fmt"
	"go/ast"
)

// Type used to represent cyclomatic complexity
// for each function.
type FunctionComplexity struct {
	FunctionName     string
	controlFlowGraph graph.Graph
}

func (functionCC *FunctionComplexity) GetCyclomaticComplexity() (int) {
	//CyclomaticComplexity = NumberOfEdges - NumberOfNodes + 2 * NumberOfConnectedComponents
	fmt.Printf("NumberOfEdges: %d\n", functionCC.controlFlowGraph.GetNumberOfEdges())
	fmt.Printf("NumberOfNodes: %d\n", functionCC.controlFlowGraph.GetNumberOfNodes())
	fmt.Printf("NumberOfConnectedComponents: %d\n", functionCC.controlFlowGraph.GetNumberOfSCComponents())
	return functionCC.controlFlowGraph.GetNumberOfEdges() - functionCC.controlFlowGraph.GetNumberOfNodes() + 2 * functionCC.controlFlowGraph.GetNumberOfSCComponents()
}

// GetCyclomaticComplexityFileLevel return the cyclomatic
// complexity for the entire file specified as argument.
func GetCyclomaticComplexityFileLevel(srcFile []byte) (int) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", srcFile, 0)
	if err != nil {
		panic(err)
	}
	controlFlowGraph := graph.NewGraph()

	for _, bb := range bblock.GetBasicBlocksFromSourceCode(file) {
		for _, sbb := range bb.Successor {
			controlFlowGraph.InsertEdge(&graph.Node{Value:bb}, &graph.Node{Value:sbb})
		}
	}
	//CyclomaticComplexity = NumberOfEdges - NumberOfNodes + 2 * NumberOfConnectedComponents
	return controlFlowGraph.GetNumberOfEdges() - controlFlowGraph.GetNumberOfNodes() + 2 * controlFlowGraph.GetNumberOfSCComponents()
}

// TODO: Should one let []byte representing the file.
func GetCyclomaticComplexityFunctionLevel(srcFile []byte) (functions []*FunctionComplexity) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", srcFile, 0)
	if err != nil {
		panic(err)
	}

	for _, bb := range bblock.GetBasicBlocksFromSourceCode(file) {
		if bb.Type == bblock.FUNCTION_ENTRY {
			functions = append(functions, &FunctionComplexity{FunctionName:bb.Head.(*ast.FuncDecl).Name.Name, controlFlowGraph:*graph.NewGraph()})
		}
		if len(functions) > 0 {
			functions[len(functions) - 1].controlFlowGraph.InsertNode(&graph.Node{Value:bb})
			for _, sbb := range bb.Successor {
				functions[len(functions) - 1].controlFlowGraph.InsertEdge(&graph.Node{Value:bb}, &graph.Node{Value:sbb})
			}
		}
	}
	return functions
}
