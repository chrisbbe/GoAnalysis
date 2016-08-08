// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package ccomplexity

import (
	"fmt"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity/bblock"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity/cfgraph"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity/graph"
)

// FunctionComplexity represents cyclomatic complexity in a function or method.
type FunctionComplexity struct {
	Name             string                    //Function name.
	SrcLine          int                       //Line number in source file where func is declared.
	Complexity       int                       //Cyclomatic complexity value.
	ControlFlowGraph *cfgraph.ControlFlowGraph //Control-flow graph in function.
	BasicBlocks      []*bblock.BasicBlock      //Basic-blocks in function.
}

func (funCC *FunctionComplexity) String() string {
	return fmt.Sprintf("%s() at line %d has CC: %d\n", funCC.Name, funCC.SrcLine, funCC.Complexity)
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

func (function *FunctionComplexity) GetSCComponents() []*graph.StronglyConnectedComponent {
	return function.ControlFlowGraph.GetSCComponents()
}

func GetCyclomaticComplexity(cfg *cfgraph.ControlFlowGraph) int {
	return cfg.GetNumberOfEdges() - cfg.GetNumberOfNodes() + cfg.GetNumberOfSCComponents()
}

func GetCyclomaticComplexityFunctionLevel(goFilePath string, goSrcFile []byte) (functions []*FunctionComplexity, err error) {
	blocks, err := bblock.GetBasicBlocksFromSourceCode(goFilePath, goSrcFile)
	if err != nil {
		return nil, err
	}

	for _, cfg := range cfgraph.GetControlFlowGraph(blocks) {
		complexity := GetCyclomaticComplexity(cfg)
		funcBlock := cfg.Root.Value.(*bblock.BasicBlock)
		functions = append(functions, &FunctionComplexity{
			Name:             funcBlock.FunctionName,
			SrcLine:          funcBlock.FunctionDeclLine,
			Complexity:       complexity,
			ControlFlowGraph: cfg,
			BasicBlocks:      blocks,
		})
	}
	return functions, nil
}
