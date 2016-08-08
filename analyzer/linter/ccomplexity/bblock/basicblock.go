// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package bblock

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"sort"
)

type BasicBlockType int

//Basic Block types.
const (
	FUNCTION_ENTRY BasicBlockType = iota
	IF_CONDITION
	ELSE_CONDITION
	SWITCH_STATEMENT
	CASE_CLAUSE
	SELECT_STATEMENT
	COMM_CLAUSE
	RETURN_STMT
	FOR_STATEMENT
	RANGE_STATEMENT
	GO_STATEMENT
	CALL_EXPRESSION
	ELSE_BODY
	FOR_BODY
	EMPTY
	START
	EXIT
	UNKNOWN
)

var basicBlockTypeStrings = [...]string{
	FUNCTION_ENTRY:   "FUNCTION_ENTRY",
	IF_CONDITION:     "IF_CONDITION",
	ELSE_CONDITION:   "ELSE_CONDITION",
	SWITCH_STATEMENT: "SWITCH_STATEMENT",
	CASE_CLAUSE:      "CASE_CLAUSE",
	SELECT_STATEMENT: "SELECT_STATEMENT",
	COMM_CLAUSE:      "COMM_CLAUSE",
	RETURN_STMT:      "RETURN_STMT",
	FOR_STATEMENT:    "FOR_STATEMENT",
	RANGE_STATEMENT:  "RANGE_STATEMENT",
	GO_STATEMENT:     "GO_STATEMENT",
	CALL_EXPRESSION:  "CALL_EXPRESSION",
	ELSE_BODY:        "ELSE_BODY",
	FOR_BODY:         "FOR_BODY",
	EMPTY:            "EMPTY",
	START:            "Start",
	EXIT:             "Exit",
	UNKNOWN:          "UNKNOWN",
}

func (bbType BasicBlockType) String() string {
	return basicBlockTypeStrings[bbType]
}

func (basicBlock *BasicBlock) UID() string {
	//Both START and EXIT blocks are meta-blocks, giving them negative UID which will not be confused with 'real' blocks.
	if basicBlock.Type == START || basicBlock.Type == EXIT {
		return fmt.Sprintf("%d", 0-basicBlock.Type)
	}
	return fmt.Sprintf("%d", basicBlock.EndLine)
}

func (basicBlock *BasicBlock) String() string {
	if basicBlock.Type == START {
		return basicBlock.Type.String()
	} else if basicBlock.Type == EXIT {
		return basicBlock.Type.String()
	}
	return fmt.Sprintf("BLOCK NR.%d (%s) (EndLine: %d)", basicBlock.Number, basicBlock.Type.String(), basicBlock.EndLine)
}

func (basicBlock *BasicBlock) AddSuccessorBlock(successorBlocks ...*BasicBlock) {
	for _, successorBlock := range successorBlocks {
		basicBlock.successor[successorBlock.EndLine] = successorBlock //Update or add.
	}
}

func NewBasicBlock(blockNumber int, blockType BasicBlockType, endLine int) *BasicBlock {
	return &BasicBlock{Number: blockNumber, Type: blockType, EndLine: endLine, successor: map[int]*BasicBlock{}}
}

func (basicBlock *BasicBlock) GetSuccessorBlocks() []*BasicBlock {
	keys := make([]int, len(basicBlock.successor))
	basicBlocks := []*BasicBlock{}

	i := 0
	for k := range basicBlock.successor {
		keys[i] = k
		i++
	}
	sort.Ints(keys) //Sort keys from map.

	//Add the basic-block into the array.
	for _, key := range keys {
		basicBlocks = append(basicBlocks, basicBlock.successor[key])
	}
	return basicBlocks
}

type BasicBlock struct {
	Number           int
	Type             BasicBlockType
	EndLine          int
	successor        map[int]*BasicBlock
	FunctionName     string
	FunctionDeclLine int
}

type visitor struct {
	basicBlocks   map[int]*BasicBlock
	sourceFileSet *token.FileSet

	lastBlock     *BasicBlock
	prevLastBlock *BasicBlock

	returnBlock  *BasicBlock
	forBlock     *BasicBlock
	forBodyBlock *BasicBlock
	switchBlock  *BasicBlock
}

// UpdateBasicBlock updates all the variables from the newBasicBlock into the basicBlock object.
func (basicBlock *BasicBlock) UpdateBasicBlock(newBasicBlock *BasicBlock) {
	if newBasicBlock != nil {
		basicBlock.Number = newBasicBlock.Number
		basicBlock.Type = newBasicBlock.Type
		basicBlock.EndLine = newBasicBlock.EndLine
		basicBlock.successor = newBasicBlock.successor
		basicBlock.FunctionName = newBasicBlock.FunctionName
		basicBlock.FunctionDeclLine = newBasicBlock.FunctionDeclLine
	}
}

func (v *visitor) AddBasicBlock(blockType BasicBlockType, position token.Pos) *BasicBlock {
	line := v.sourceFileSet.File(position).Line(position)
	basicBlock := NewBasicBlock(-1, blockType, line) //-1 indicates number will be set later.

	// Bookkeeping.
	v.prevLastBlock = v.lastBlock
	v.lastBlock = basicBlock

	//Update the existing block., or add new block.
	if bb, ok := v.basicBlocks[line]; ok {
		bb.UpdateBasicBlock(basicBlock)
		v.lastBlock = bb
		return bb
	} else {
		v.basicBlocks[line] = basicBlock
	}
	return basicBlock
}

// GetBasicBlocks converts map holding the basic-blocks to the ordered set
// of basic-blocks, in right order!
func (v *visitor) GetBasicBlocks() []*BasicBlock {
	keys := make([]int, len(v.basicBlocks))
	basicBlocks := make([]*BasicBlock, len(v.basicBlocks))

	i := 0
	for k := range v.basicBlocks {
		keys[i] = k
		i++
	}
	sort.Ints(keys) //Sort keys from map.

	//Add the basic-block into the array.
	for index, key := range keys {
		basicBlocks[index] = v.basicBlocks[key]
		basicBlocks[index].Number = index //Set basic-block number.
	}
	return basicBlocks
}

func GetBasicBlocksFromSourceCode(filePath string, srcFile []byte) ([]*BasicBlock, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filePath, srcFile, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("Parse error: %s", err)
	}

	visitor := &visitor{sourceFileSet: fileSet, basicBlocks: make(map[int]*BasicBlock)}
	ast.Walk(visitor, file)

	basicBlocks := visitor.GetBasicBlocks()

	numberOfBasicBlocks := len(basicBlocks)
	for index, bBlock := range basicBlocks {
		if bBlock.Type != FOR_BODY && bBlock.Type != ELSE_CONDITION && bBlock.Type != ELSE_BODY && bBlock.Type != COMM_CLAUSE && bBlock.Type != CASE_CLAUSE && bBlock.Type != RETURN_STMT {
			if numberOfBasicBlocks > index+1 {
				bBlock.AddSuccessorBlock(basicBlocks[index+1])
			}
		}
	}

	return basicBlocks, nil
}

func PrintBasicBlocks(basicBlocks []*BasicBlock) {
	for _, bb := range basicBlocks {
		log.Printf("%d) %s (EndLine: %d)\n", bb.Number, bb.Type.String(), bb.EndLine)

		for _, sBB := range bb.GetSuccessorBlocks() {
			log.Printf("\t-> (%d) %s (EndLine: %d)\n", sBB.Number, sBB.Type.String(), sBB.EndLine)
		}
	}
}

//TODO: Check after all basic-block types we have declared.
func GetBasicBlockTypeFromStmt(stmtList []ast.Stmt) (BasicBlockType, ast.Stmt) {
	for _, stmt := range stmtList {
		switch stmt.(type) {
		case *ast.ReturnStmt:
			return RETURN_STMT, stmt
		case *ast.CaseClause:
			return CASE_CLAUSE, stmt
		case *ast.SwitchStmt:
			return SWITCH_STATEMENT, stmt
		}
	}
	return UNKNOWN, nil
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch t := node.(type) {

		case *ast.FuncDecl:
			funcDeclBlock := v.AddBasicBlock(FUNCTION_ENTRY, t.Pos())
			// Set information about Function in the block.
			funcDeclBlock.FunctionName = t.Name.Name
			funcDeclBlock.FunctionDeclLine = v.sourceFileSet.File(t.Pos()).Line(t.Pos())

			if t.Body != nil {
				for _, s := range t.Body.List {
					if _, ok := s.(*ast.ReturnStmt); ok {
						v.returnBlock = v.AddBasicBlock(RETURN_STMT, s.End())
					}
				}

				if v.returnBlock == nil {
					v.returnBlock = v.AddBasicBlock(RETURN_STMT, t.End())
				}

				//Visit all statements in body.
				for _, s := range t.Body.List {
					v.Visit(s)
				}
			}

			v.returnBlock = nil
			return nil

		case *ast.ReturnStmt:
			prevBlock := v.lastBlock
			v.returnBlock = v.AddBasicBlock(RETURN_STMT, t.Pos())

			if prevBlock != nil && prevBlock.Type != RETURN_STMT && len(prevBlock.successor) == 0 {
				prevBlock.AddSuccessorBlock(v.returnBlock)
			}

			if v.switchBlock != nil {
				v.switchBlock.AddSuccessorBlock(v.returnBlock)
			}

		case *ast.GoStmt:
			v.AddBasicBlock(GO_STATEMENT, t.Pos())

		case *ast.IfStmt:
			ifBlock := v.AddBasicBlock(IF_CONDITION, t.Pos())

			for _, stmt := range t.Body.List {
				v.Visit(stmt)
			}

			var elseConditionBlock, elseBodyBlock *BasicBlock
			if t.Else != nil {

				if v.lastBlock != nil && v.lastBlock.Type != RETURN_STMT {
					elseConditionBlock = v.AddBasicBlock(ELSE_CONDITION, t.Else.Pos())
				} else {
					// We don't want to set return block as successor to elseCondition.
					v.returnBlock = nil
				}

				elseBodyBlock = v.AddBasicBlock(ELSE_BODY, t.Else.End())
				ifBlock.AddSuccessorBlock(elseBodyBlock)

				if v.returnBlock != nil {
					if elseConditionBlock != nil {
						elseConditionBlock.AddSuccessorBlock(v.returnBlock)
					}
					elseBodyBlock.AddSuccessorBlock(v.returnBlock)
				}
			}

		case *ast.ForStmt:
			v.forBlock = v.AddBasicBlock(FOR_STATEMENT, t.Pos())

			if v.returnBlock != nil {
				v.forBlock.AddSuccessorBlock(v.returnBlock)
			}

			tmpReturnBlock := v.returnBlock
			tmpForBlock := v.forBlock
			v.returnBlock = v.forBlock

			for _, s := range t.Body.List {
				v.Visit(s)
			}

			v.returnBlock = tmpReturnBlock
			v.forBlock = tmpForBlock

			if v.lastBlock.Type == FOR_STATEMENT {
				v.AddBasicBlock(FOR_BODY, t.End())
			}

			if v.lastBlock.Type != RETURN_STMT {
				v.lastBlock.AddSuccessorBlock(v.forBlock)
			}

			v.forBlock = nil
			return nil

		case *ast.SwitchStmt:
			v.switchBlock = v.AddBasicBlock(SWITCH_STATEMENT, t.Pos())
			if v.forBlock != nil {
				v.forBlock.AddSuccessorBlock(v.switchBlock)
				v.switchBlock.AddSuccessorBlock(v.forBlock)
			}

			if v.returnBlock != nil {
				v.switchBlock.AddSuccessorBlock(v.returnBlock)
			}

			for _, s := range t.Body.List {
				v.Visit(s)
			}
			return nil

		case *ast.TypeSwitchStmt:
			v.switchBlock = v.AddBasicBlock(SWITCH_STATEMENT, t.Pos())
			if v.forBlock != nil {
				v.forBlock.AddSuccessorBlock(v.switchBlock)
				v.switchBlock.AddSuccessorBlock(v.forBlock)
			}

			for _, s := range t.Body.List {
				v.Visit(s)
			}
			return nil

		case *ast.SelectStmt:
			v.switchBlock = v.AddBasicBlock(SELECT_STATEMENT, t.Pos())
			if v.forBlock != nil {
				v.forBlock.AddSuccessorBlock(v.switchBlock)
				v.switchBlock.AddSuccessorBlock(v.forBlock)
			}

			for _, s := range t.Body.List {
				v.Visit(s)
			}
			return nil

		case *ast.CaseClause:
			var caseClause *BasicBlock
			if basicBlockType, s := GetBasicBlockTypeFromStmt(t.Body); basicBlockType != UNKNOWN {
				caseClause = v.AddBasicBlock(basicBlockType, s.Pos())
			} else {
				caseClause = v.AddBasicBlock(CASE_CLAUSE, t.End())
			}

			if v.forBlock != nil {
				caseClause.AddSuccessorBlock(v.forBlock)
			}

			if v.switchBlock != nil {
				v.switchBlock.AddSuccessorBlock(caseClause)
			}

			if v.returnBlock != nil {
				caseClause.AddSuccessorBlock(v.returnBlock)
			}

			tmpSwitchBlock := v.switchBlock
			tmpReturnBLock := v.returnBlock
			for _, s := range t.Body {
				v.Visit(s)
			}
			v.switchBlock = tmpSwitchBlock
			v.returnBlock = tmpReturnBLock

			//TODO: Special case.
			//TODO: Type is always CASE_CLAUSE type
			if v.returnBlock != nil && caseClause.Type != RETURN_STMT && caseClause.Type != SWITCH_STATEMENT {
				//TODO: This must be refactored more beautiful
				containsForStatement := false
				for _, b := range caseClause.GetSuccessorBlocks() {
					if b.Type == FOR_STATEMENT {
						containsForStatement = true
					}
				}
				if !containsForStatement {
					caseClause.AddSuccessorBlock(v.returnBlock)
				}
			}

		case *ast.CommClause:
			var caseClause *BasicBlock
			if basicBlockType, s := GetBasicBlockTypeFromStmt(t.Body); basicBlockType != UNKNOWN {
				caseClause = v.AddBasicBlock(basicBlockType, s.Pos())
			} else {
				caseClause = v.AddBasicBlock(COMM_CLAUSE, t.End())
			}

			if v.forBlock != nil {
				caseClause.AddSuccessorBlock(v.forBlock)
			}

			if v.switchBlock != nil {
				v.switchBlock.AddSuccessorBlock(caseClause)
			}

			if v.returnBlock != nil {
				caseClause.AddSuccessorBlock(v.returnBlock)
			}

			tmpSwitchBlock := v.switchBlock
			tmpReturnBLock := v.returnBlock
			for _, s := range t.Body {
				v.Visit(s)
			}
			v.switchBlock = tmpSwitchBlock
			v.returnBlock = tmpReturnBLock

			//TODO: Special case.
			//TODO: Type is always CASE_CLAUSE type
			if v.returnBlock != nil && caseClause.Type != RETURN_STMT && caseClause.Type != SWITCH_STATEMENT {
				//TODO: This must be refactored more beautifully
				containsForStatement := false
				for _, b := range caseClause.GetSuccessorBlocks() {
					if b.Type == FOR_STATEMENT {
						containsForStatement = true
					}
				}
				if !containsForStatement {
					caseClause.AddSuccessorBlock(v.returnBlock)
				}
			}
		}
	}
	return v
}
