package bblock

import (
	"go/ast"
	"go/token"
)

type basicBlock struct {
	Number   int
	FromLine int
	ToLine   int
}

type visitor struct {
	fileSet     *token.FileSet
	basicBlocks []*basicBlock
	leader      *basicBlock
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch t := node.(type) {

		case *ast.IfStmt:
			v.leader.ToLine = getSourceCodeLine(v.fileSet, t.Pos()) - 1 //Update previous leader with toLine
			bb := basicBlock{Number: len(v.basicBlocks), FromLine: getSourceCodeLine(v.fileSet, t.Pos())}
			v.basicBlocks = append(v.basicBlocks, &bb)
			v.leader = &bb

		case *ast.SwitchStmt:
			v.leader.ToLine = getSourceCodeLine(v.fileSet, t.Pos()) - 1  //Update previous leader with toLine
			bb := basicBlock{Number: len(v.basicBlocks), FromLine: getSourceCodeLine(v.fileSet, t.Pos())}
			v.basicBlocks = append(v.basicBlocks, &bb)
			v.leader = &bb

		case *ast.CaseClause:
			v.leader.ToLine = getSourceCodeLine(v.fileSet, t.Pos()) - 1  //Update previous leader with toLine
			bb := basicBlock{Number: len(v.basicBlocks), FromLine: getSourceCodeLine(v.fileSet, t.Pos())}
			v.basicBlocks = append(v.basicBlocks, &bb)
			v.leader = &bb

		default:
			if len(v.basicBlocks) == 0 { //Create start BBlock
				bb := basicBlock{Number: len(v.basicBlocks), FromLine: getSourceCodeLine(v.fileSet, t.Pos())}
				v.basicBlocks = append(v.basicBlocks, &bb)
				v.leader = &bb
			} else { //Update last leaders blocks toLine.
				v.leader.ToLine = getSourceCodeLine(v.fileSet, t.Pos())
			}
		}
	}
	return v
}

func GetBasicBlocksFromSourceCode(fset *token.FileSet, file *ast.File) ([]*basicBlock) {
	fv := visitor{fileSet: fset}
	ast.Walk(&fv, file)
	return fv.basicBlocks
}

func getSourceCodeLine(fileSet *token.FileSet, pos token.Pos) int {
	return fileSet.File(pos).Line(pos)
}
