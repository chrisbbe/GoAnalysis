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
	"go/ast"
)

type BasicBlockType int

//Basic Block types.
const (
	PACKAGE_ENTRY = iota //0
	FUNCTION_ENTRY = iota //1
	IF_CONDITION = iota //2
	ELSE_CONDITION = iota //3
	SWITCH_CONDITION = iota //4
	CASE_CLAUSE = iota //5
	SELECT_STATEMENT = iota // 6
	COMM_CLAUSE = iota // 7
	RETURN_STMT = iota //8
)

func (bbType BasicBlockType) String() string {
	switch bbType {
	case PACKAGE_ENTRY:
		return "PACKAGE_ENTRY"
	case FUNCTION_ENTRY:
		return "FUNCTION_ENTRY"
	case IF_CONDITION:
		return "IF_CONDITION"
	case ELSE_CONDITION:
		return "ELSE_CONDITION"
	case SWITCH_CONDITION:
		//TODO: Should we name it switch_condition or switch_statement ?
		return "SWITCH_CONDITION"
	case CASE_CLAUSE:
		return "CASE_CLAUSE"
	case SELECT_STATEMENT:
		return "SELECT_STATEMENT"
	case COMM_CLAUSE:
		return "COMM_CLAUSE"
	case RETURN_STMT:
		return "RETURN_STMT"
	}
	return "Unknown"
}

type SwitchBlock struct {
	index int
	Cases []ast.Stmt
}

type SelectBlock struct {
	index       int
	CommClauses []ast.Stmt
}

type IfBlock struct {
	index     int
	ElseStmt  *ast.Stmt
	ElseIndex int
}

type visitor struct {
	basicBlocks    []*BasicBlock
	SwitchIndexes  []*SwitchBlock
	SelectIndexes  []*SelectBlock
	IfIndexes      []*IfBlock //TODO: Do we actually need this one?
	currentAstNode *ast.Node
	prevAstNode    *ast.Node
}

type BasicBlock struct {
	Number    int            //Successive incremental block number.
	Type      BasicBlockType //Type corresponds to the type starting the BB.
	Head      ast.Node       //First Node in BB.
	Tail      ast.Node       //Last Node in BB.
	Successor []*BasicBlock  //Successor BB.
	IfBlock   *BasicBlock    //Used to point on If-block in else.
}

func (v *visitor) bindElseBlocks() {
	for i := len(v.basicBlocks) - 1; i >= 0; i-- {
		bb := v.basicBlocks[i]

		if bb.Type == IF_CONDITION {
			v.basicBlocks[i - 1].AddSuccessorBlock(bb)
		} else if bb.Type == ELSE_CONDITION {
			v.basicBlocks[bb.IfBlock.Number - 1].AddSuccessorBlock(bb)
		}
	}
}

func (v *visitor) AddBasicBlock(blockType BasicBlockType, nodeType ast.Node) {
	if v.prevAstNode != nil {
		// Update tail on previous block when we are creating a new one!.
		v.basicBlocks[len(v.basicBlocks) - 1].Tail = *v.prevAstNode
	}
	v.basicBlocks = append(v.basicBlocks, &BasicBlock{Type: blockType, Number: len(v.basicBlocks), Head: nodeType})
}

func (basicBlock *BasicBlock) AddSuccessorBlock(bb *BasicBlock) {
	basicBlock.Successor = append(basicBlock.Successor, bb)
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		//Update pointers to previous and current node.
		v.prevAstNode = v.currentAstNode
		v.currentAstNode = &node

		switch astNodeType := node.(type) {

		case *ast.File:
			v.AddBasicBlock(PACKAGE_ENTRY, astNodeType)

		case *ast.FuncDecl:
			v.AddBasicBlock(FUNCTION_ENTRY, astNodeType)

		case *ast.GoStmt:
			//TODO: Is Go stmt the same as function entry?
			v.AddBasicBlock(FUNCTION_ENTRY, astNodeType)

		//TODO: Unsure if RETURN should generate a new block.
		//case *ast.ReturnStmt:
		//v.AddBasicBlock(RETURN_STMT, astNodeType)

		case *ast.IfStmt:
			v.AddBasicBlock(IF_CONDITION, astNodeType)

			if astNodeType.Else != nil {
				//If has Else part.
				v.IfIndexes = append(v.IfIndexes, &IfBlock{index:len(v.basicBlocks) - 1, ElseStmt: &astNodeType.Else})
			}

		case *ast.SwitchStmt:
			v.AddBasicBlock(SWITCH_CONDITION, astNodeType)
			v.SwitchIndexes = append(v.SwitchIndexes, &SwitchBlock{index: len(v.basicBlocks) - 1, Cases: astNodeType.Body.List})

		case *ast.TypeSwitchStmt:
			v.AddBasicBlock(SWITCH_CONDITION, astNodeType)
			v.SwitchIndexes = append(v.SwitchIndexes, &SwitchBlock{index: len(v.basicBlocks) - 1, Cases: astNodeType.Body.List})

		case *ast.CaseClause:
			v.AddBasicBlock(CASE_CLAUSE, astNodeType)
			for _, sbb := range v.SwitchIndexes {
				for _, bb := range sbb.Cases {
					if astNodeType == bb {
						v.basicBlocks[sbb.index].AddSuccessorBlock(v.basicBlocks[len(v.basicBlocks) - 1])
					}
				}
			}

		case *ast.SelectStmt:
			v.AddBasicBlock(SELECT_STATEMENT, astNodeType)
			v.SelectIndexes = append(v.SelectIndexes, &SelectBlock{index: len(v.basicBlocks) - 1, CommClauses:astNodeType.Body.List})

		case *ast.CommClause:
			v.AddBasicBlock(COMM_CLAUSE, astNodeType)
			for _, sbb := range v.SelectIndexes {
				for _, bb := range sbb.CommClauses {
					if astNodeType == bb {
						v.basicBlocks[sbb.index].AddSuccessorBlock(v.basicBlocks[len(v.basicBlocks) - 1])
					}
				}
			}

		case ast.Stmt:
			//Check if 'astNodeType' equals an Else-stmt.
			for _, ifBlock := range v.IfIndexes {
				if astNodeType == *ifBlock.ElseStmt {
					v.AddBasicBlock(ELSE_CONDITION, astNodeType)
					v.basicBlocks[len(v.basicBlocks) - 1].IfBlock = v.basicBlocks[ifBlock.index]
				}
			}

		default:
			if len(v.basicBlocks) > 0 {
				//Updating last BB EndNode.
				v.basicBlocks[len(v.basicBlocks) - 1].Tail = astNodeType
			}
		}
	}
	return v
}

func GetBasicBlocksFromSourceCode(file *ast.File) []*BasicBlock {
	fv := visitor{}
	ast.Walk(&fv, file)

	fv.bindElseBlocks()

	// Set basic-block 'e' as successor for the previous basic-block.
	for index, bb := range fv.basicBlocks {
		if len(bb.Successor) == 0 && len(fv.basicBlocks) > index + 1 {
			if bb.Type != CASE_CLAUSE && bb.Type != COMM_CLAUSE && bb.Type != IF_CONDITION && bb.Type != ELSE_CONDITION {
				bb.AddSuccessorBlock(fv.basicBlocks[index + 1])
			}
		}
	}
	return fv.basicBlocks
}
