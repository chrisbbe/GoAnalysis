package bblock

import (
	"go/ast"
	"go/token"
	"fmt"
)

type basicBlock struct {
	Number   int
	FromLine int
	ToLine   int
	Value    string
}

type visitor struct {
	fileSet     *token.FileSet
	basicBlocks []*basicBlock
	leader      *basicBlock
	preGoTo     bool
	ifEnd       token.Pos
	elseEnd     token.Pos
}

func (v *visitor) addNewBasicBlock(node ast.Node, preGoTo bool, value string) {
	v.leader.ToLine = getLineInSourceCode(v.fileSet, node.Pos()) - 1 //Update previous leader with toLine
	bb := basicBlock{Value: value, Number: len(v.basicBlocks), FromLine: getLineInSourceCode(v.fileSet, node.Pos())}
	v.basicBlocks = append(v.basicBlocks, &bb)
	v.leader = &bb
	v.preGoTo = preGoTo
}

//TODO: Is it necessarily with the parameter?
func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch t := node.(type) {

		case *ast.IfStmt: //Detecting both if and else-if.
			v.addNewBasicBlock(t, true, "If")
			v.ifEnd = t.End()
			if t.Else != nil {
				v.elseEnd = t.Else.End()
			}

		case *ast.SwitchStmt:
			v.addNewBasicBlock(t, true, "Switch")

		case *ast.CaseClause:
			v.addNewBasicBlock(t, true, "CaseClause")

		default:
			if len(v.basicBlocks) == 0 { //Create start BBlock
				bb := basicBlock{Value: "First", Number: len(v.basicBlocks), FromLine: getLineInSourceCode(v.fileSet, t.Pos())}
				v.basicBlocks = append(v.basicBlocks, &bb)
				v.leader = &bb
			} else { //Update last leaders blocks toLine.
				v.leader.ToLine = getLineInSourceCode(v.fileSet, t.Pos())
				if v.preGoTo {
					//v.leader.ToLine = getLineInSourceCode(v.fileSet, t.Pos()) - 1  //Update previous leader with toLine
					bb := basicBlock{Value: "After Go", Number: len(v.basicBlocks), FromLine: getLineInSourceCode(v.fileSet, t.Pos())}
					v.basicBlocks = append(v.basicBlocks, &bb)
					v.leader = &bb
					v.preGoTo = false
				}
			}

			if v.ifEnd.IsValid() && t.End() >= v.ifEnd {
				fmt.Println("Ifdone\n new block")
				v.addNewBasicBlock(t, true, "IfDone")
				v.ifEnd = token.NoPos
			}

			//Maybe call isValid() to check if position is valid
			if v.elseEnd.IsValid() && t.End() >= v.elseEnd {
				fmt.Println("Else done\n new block")
				v.addNewBasicBlock(t, true, "ElseDone")
				v.elseEnd = token.NoPos
			}
		}
		fmt.Printf("Node = %T\n", node)
	}
	return v
}

//TODO: What is fset?? Rename to better name!
func GetBasicBlocksFromSourceCode(fset *token.FileSet, file *ast.File) ([]*basicBlock) {
	//TODO: According to the documentation, we should use: NoPos to indicate no-postion.
	fv := visitor{fileSet: fset, preGoTo:false}
	ast.Walk(&fv, file)
	return fv.basicBlocks
}

func getLineInSourceCode(fileSet *token.FileSet, pos token.Pos) int {
	return fileSet.File(pos).Line(pos)
}
