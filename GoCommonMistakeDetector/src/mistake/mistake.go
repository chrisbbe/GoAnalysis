package mistake

import (
	"go/ast"
	"go/token"
	"fmt"
	"go/parser"
	"errors"
)

//Mistake types.
const (
	RACE_CONDITION = iota  // 0
	FMT_PRINTING = iota // 1
)

//Representing mistakes found when
//analysis source-code.
type Mistake struct {
	LineInSourceFile int
	MistakeType      int
}

//Type used by Visitor to
//hold data when walking
//the tree. mistakes is
//a slice holding all 'mistakes'
//found.
type visitor struct {
	fileSet  *token.FileSet
	mistakes []Mistake
}

//FindCommonMistakes is the function creating the AST and
//initiating walking of the tree, and returning final result
// of analysis.
func FindCommonMistakes(sourceFile []byte) (mistakes []Mistake, walkError error) {
	// Create the AST by parsing src.
	fileSet := token.NewFileSet() // positions are relative to fileSet
	f, err := parser.ParseFile(fileSet, "", sourceFile, 0)
	if err != nil {
		return nil, err
	}

	//Visit() method may panic(),
	//which is catched in this deferred
	//closure.
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				walkError = errors.New(x)
			case error:
				walkError = x
			default:
				walkError = errors.New("Panic unknown")
			}
		}
	}()

	visitor := visitor{fileSet: fileSet}
	ast.Walk(&visitor, f)
	return visitor.mistakes, walkError
}

func (visitor *visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.GoStmt:
		if findRaceInGoRoutine(t) {
			sourceLineNumber := getLineNumberInSourceCode(visitor.fileSet, t.Pos())
			newMistake := Mistake{LineInSourceFile:sourceLineNumber, MistakeType:RACE_CONDITION}
			visitor.mistakes = append(visitor.mistakes, newMistake)
		}
	case *ast.SelectorExpr:
		value, ok := t.X.(*ast.Ident)
		if ok && value.Name == "fmt" && (t.Sel.Name == "Print" || t.Sel.Name == "Println" || t.Sel.Name == "Printf" || t.Sel.Name == "Fprintln" || t.Sel.Name == "Fprintf" || t.Sel.Name == "Errorf") {
			sourceLineNumber := getLineNumberInSourceCode(visitor.fileSet, t.Pos())
			newMistake := Mistake{LineInSourceFile:sourceLineNumber, MistakeType:FMT_PRINTING}
			visitor.mistakes = append(visitor.mistakes, newMistake)
		}

	//TODO: Should we care about everything that is bad, or is it catched somwhere else?
	case *ast.BadDecl:
		sourceLineNumber := getLineNumberInSourceCode(visitor.fileSet, t.Pos())
		panic(fmt.Sprintf("Error: Parse error at line %d.\n", sourceLineNumber))
	case *ast.BadExpr:
		sourceLineNumber := getLineNumberInSourceCode(visitor.fileSet, t.Pos())
		panic(fmt.Sprintf("Error: Parse error at line %d.\n", sourceLineNumber))
	case *ast.BadStmt:
		sourceLineNumber := getLineNumberInSourceCode(visitor.fileSet, t.Pos())
		panic(fmt.Sprintf("Error: Parse error at line %d.\n", sourceLineNumber))
	}
	return visitor
}

func findRaceInGoRoutine(goNode *ast.GoStmt) (races bool) {
	switch t := goNode.Call.Fun.(type) {
	case *ast.FuncLit:
		params := t.Type.Params.List
		for _, each := range t.Body.List {
			switch t1 := each.(type) {
			case *ast.ExprStmt:
				if !validateParams(t1, params) {
					return true
				}
			}
		}
	}
	return false
}

func validateParams(node *ast.ExprStmt, List []*ast.Field) (valid bool) {
	switch t := node.X.(type) {
	case *ast.CallExpr:
		for _, each := range t.Args {
			switch t1 := each.(type) {
			case *ast.Ident:
				if !containsListParam(t1, List) {
					return false
				}
			}
		}
	}
	return true
}

func containsListParam(ident *ast.Ident, List []*ast.Field) (found bool) {
	for _, each := range List {
		for _, each1 := range each.Names {
			if each1.Name == ident.Name {
				return true
			}
		}
	}
	return false
}

func getLineNumberInSourceCode(fileSet *token.FileSet, position token.Pos) (line int) {
	tokenFile := fileSet.File(position)
	return tokenFile.Line(position)
}
