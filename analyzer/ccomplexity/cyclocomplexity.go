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
// of Go source code files, both cyclomatic complexity per function and file level
// is supported.
package ccomplexity

import (
	"github.com/chrisbbe/GoAnalysis/analyzer/ccomplexity/bblock"
	"github.com/chrisbbe/GoAnalysis/analyzer/ccomplexity/cfgraph"
)

// FunctionComplexity represents cyclomatic complexity in a function or method.
type FunctionComplexity struct {
	Name             string                    //Function name.
	Line             int                       //Line number in source file.
	Complexity       int                       //Cyclomatic complexity value.
	ControlFlowGraph *cfgraph.ControlFlowGraph //Control flow ccomplexity.graph in function.
	BasicBlocks      []*bblock.BasicBlock      //Basic-blocks in function.
}

func (function *FunctionComplexity) GetNumberOfNodes() int {
	return function.ControlFlowGraph.GetNumberOfNodes()
}

func (function *FunctionComplexity) GetNumberOfEdges() int {
	return function.ControlFlowGraph.GetNumberOfEdges()
}

func (function *FunctionComplexity) GetNumberOfSCC() int {
	return function.ControlFlowGraph.GetNumberOfSCComponents()
}

func GetCyclomaticComplexity(cfg *cfgraph.ControlFlowGraph) int {
	return cfg.GetNumberOfEdges() - cfg.GetNumberOfNodes() + cfg.GetNumberOfSCComponents()
}

func GetCyclomaticComplexityFunctionLevel(srcFile []byte) (functions []*FunctionComplexity, err error) {
	blocks, err := bblock.GetBasicBlocksFromSourceCode(srcFile)
	if err != nil {
		return nil, err
	}

	for _, cfg := range cfgraph.GetControlFlowGraph(blocks) {
		complexity := GetCyclomaticComplexity(cfg)
		functions = append(functions, &FunctionComplexity{Name: cfg.Root.Value.(*bblock.BasicBlock).FunctionName, Line: blocks[0].EndLine, Complexity: complexity, ControlFlowGraph: cfg, BasicBlocks: blocks})
	}
	return functions, nil
}
