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
	"go/token"
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
	RETURN_STMT = iota
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
		return "SWITCH_CONDITION"
	case CASE_CLAUSE:
		return "CASE_CLAUSE"
	case RETURN_STMT:
		return "RETURN_STMT"
	}
	return "Unknown"
}

type SwitchBlock struct {
	index int
	Cases []ast.Stmt
}

type IfBlock struct {
	index int
	used  bool
}

type visitor struct {
	basicBlocks    []*BasicBlock
	ElseStmts      []ast.Stmt
	SwitchIndexes  []*SwitchBlock
	IfStmtIndexes  []*IfBlock

	currentAstNode *ast.Node
	prevAstNode    *ast.Node
}

type BasicBlock struct {
	Number    int
	Type      BasicBlockType //Type corresponds to the type starting the BB.
	Head      ast.Node       //First Node in BB.
	Tail      ast.Node       //Last Node in BB.
	Successor []*BasicBlock  //Successor BB.
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		//Update pointers to prev and current node.
		v.prevAstNode = v.currentAstNode
		v.currentAstNode = &node

		switch t := node.(type) {

		case *ast.File:
			v.basicBlocks = append(v.basicBlocks, &BasicBlock{Type:PACKAGE_ENTRY, Number:len(v.basicBlocks), Head:t})

		case *ast.FuncDecl:
			v.updateTailOnPrevBasicBlock(v.prevAstNode)
			v.basicBlocks = append(v.basicBlocks, &BasicBlock{Type:FUNCTION_ENTRY, Number:len(v.basicBlocks), Head:t})

		case *ast.ReturnStmt:
			v.updateTailOnPrevBasicBlock(v.prevAstNode)
			v.basicBlocks = append(v.basicBlocks, &BasicBlock{Type:RETURN_STMT, Number:len(v.basicBlocks), Head:t})

		case *ast.IfStmt:
			v.updateTailOnPrevBasicBlock(v.prevAstNode)
			v.basicBlocks = append(v.basicBlocks, &BasicBlock{Type:IF_CONDITION, Number:len(v.basicBlocks), Head:t})
			v.ElseStmts = append(v.ElseStmts, t.Else)

			//Add indexes to handle nested if-else.
			v.IfStmtIndexes = append(v.IfStmtIndexes, &IfBlock{index:len(v.basicBlocks) - 1, used:false})

		case *ast.CaseClause:
			v.updateTailOnPrevBasicBlock(v.prevAstNode)
			v.basicBlocks = append(v.basicBlocks, &BasicBlock{Type: CASE_CLAUSE, Number:len(v.basicBlocks), Head:t})

			for _, sbb := range v.SwitchIndexes {
				for _, bb := range sbb.Cases {
					if t == bb {
						v.basicBlocks[sbb.index].Successor = append(v.basicBlocks[sbb.index].Successor, v.basicBlocks[len(v.basicBlocks) - 1])
					}
				}
			}

		case *ast.SwitchStmt:
			v.updateTailOnPrevBasicBlock(v.prevAstNode)
			v.basicBlocks = append(v.basicBlocks, &BasicBlock{Type:SWITCH_CONDITION, Number:len(v.basicBlocks), Head:t})

			v.SwitchIndexes = append(v.SwitchIndexes, &SwitchBlock{index:len(v.basicBlocks) - 1, Cases:t.Body.List})

		case *ast.TypeSwitchStmt:
			v.updateTailOnPrevBasicBlock(v.prevAstNode)
			v.basicBlocks = append(v.basicBlocks, &BasicBlock{Type:SWITCH_CONDITION, Number:len(v.basicBlocks), Head:t})

		case ast.Stmt:
			for i, s := range v.ElseStmts { //Check if statement is Else-stmt.
				if s == t {
					v.updateTailOnPrevBasicBlock(v.prevAstNode)
					v.basicBlocks = append(v.basicBlocks, &BasicBlock{Type:ELSE_CONDITION, Number:len(v.basicBlocks), Head:t})
					v.ElseStmts = append(v.ElseStmts[:i], v.ElseStmts[i + 1:]...) //Remove Else-statement from slice!

					for i := len(v.IfStmtIndexes) - 1; i >= 0; i-- {
						if !v.IfStmtIndexes[i].used {
							if v.basicBlocks[len(v.basicBlocks) - 1].Type != IF_CONDITION {
								v.basicBlocks[v.IfStmtIndexes[i].index].Successor = append(v.basicBlocks[v.IfStmtIndexes[i].index].Successor, v.basicBlocks[len(v.basicBlocks) - 1])
							}
							v.IfStmtIndexes[i].used = true
						}
					}

				}
			}

		default:
			if len(v.basicBlocks) > 0 { //Updating last BB EndNode.
				v.basicBlocks[len(v.basicBlocks) - 1].Tail = t
			}
		}
	}
	return v
}

//Updates tail on current basic-block, called before new basic block is created.
func (v *visitor) updateTailOnPrevBasicBlock(prevNode *ast.Node) {
	if prevNode != nil {
		v.basicBlocks[len(v.basicBlocks) - 1].Tail = *prevNode
	}
}

func getLineInSourceCode(fileSet *token.FileSet, pos ast.Node) int {
	if pos == nil {
		return -1
	}
	return fileSet.File(pos.Pos()).Line(pos.Pos())
}

func GetBasicBlocksFromSourceCode(file *ast.File) ([]*BasicBlock) {
	fv := visitor{}
	ast.Walk(&fv, file)

	// Set basic-block 'e' as successor for the previous basic-block.
	for index, bb := range fv.basicBlocks {
		if len(bb.Successor) == 0 && len(fv.basicBlocks) > index + 1 {
			if bb.Type != CASE_CLAUSE && bb.Type != IF_CONDITION && bb.Type != ELSE_CONDITION {
				bb.Successor = append(bb.Successor, fv.basicBlocks[index + 1])
			}
		}
	}

	return fv.basicBlocks
}

