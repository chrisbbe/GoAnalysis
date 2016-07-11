// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
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
	ControlFlowGraph *cfgraph.ControlFlowGraph //Control-flow graph in function.
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
		functions = append(functions, &FunctionComplexity{
			Name:             cfg.Root.Value.(*bblock.BasicBlock).FunctionName,
			Line:             blocks[0].EndLine,
			Complexity:       complexity,
			ControlFlowGraph: cfg,
			BasicBlocks:      blocks,
		})
	}
	return functions, nil
}
