// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package linter

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"regexp"
	"runtime/debug"
)

type Rule int

//Rules.
const (
	RACE_CONDITION Rule = iota
	FMT_PRINTING
	STRING_DEFINES_ITSELF
	MAP_ALLOCATED_WITH_NEW
	ERROR_IGNORED
	EMPTY_IF_BODY
	EMPTY_ELSE_BODY
	EMPTY_FOR_BODY
	GOTO_USED
	CONDITION_EVALUATED_STATICALLY
	BUFFER_NOT_FLUSHED
	EARLY_RETURN
	NO_ELSE_RETURN
)

var ruleStrings = map[Rule]string{
	RACE_CONDITION:                 "RACE_CONDITION",
	FMT_PRINTING:                   "FMT_PRINTING",
	STRING_DEFINES_ITSELF:          "STRING_DEFINES_ITSELF",
	MAP_ALLOCATED_WITH_NEW:         "MAP_ALLOCATED_WITH_NEW",
	ERROR_IGNORED:                  "ERROR_IGNORED",
	EMPTY_IF_BODY:                  "EMPTY_IF_BODY",
	EMPTY_ELSE_BODY:                "EMPTY_ELSE_BODY",
	EMPTY_FOR_BODY:                 "EMPTY_FOR_BODY",
	GOTO_USED:                      "GOTO_USED",
	CONDITION_EVALUATED_STATICALLY: "CONDITION_EVALUATED_STATICALLY",
	BUFFER_NOT_FLUSHED:             "NO_BUFFERED_FLUSHING",
	EARLY_RETURN:                   "EARLY_RETURN",
	NO_ELSE_RETURN:                 "NO_ELSE_RETURN",
}

func (rule Rule) String() string {
	return ruleStrings[rule]
}

func DetectViolations(goSrcFile string) (violations []*Violation, err error) {
	srcFile, err := ioutil.ReadFile(goSrcFile)
	if err != nil {
		return violations, err
	}

	// Create the AST by parsing src.
	fileSet := token.NewFileSet() // positions are relative to fileSet
	fi, err := parser.ParseFile(fileSet, "", srcFile, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return violations, err
	}

	// Visit() method may panic()
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			switch x := r.(type) {
			case string:
				fmt.Print(errors.New(x))
			case error:
				fmt.Print(x)
			default:
				fmt.Print(errors.New("Panic unknown"))
			}
		}
	}()

	goFile := file{
		goFile:   fi,
		fileSet:  fileSet,
		filePath: goSrcFile,
	}

	conf := types.Config{Importer: importer.Default()}
	goFile.typeInfo = &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	//TODO: What should pathstring be ?
	if _, err := conf.Check("cmd/hello", goFile.fileSet, []*ast.File{goFile.goFile}, goFile.typeInfo); err != nil {
		return violations, err
	}

	goFile.Analyse()

	return goFile.violations, nil

}

func (goFile *file) Analyse() {
	goFile.detectFmtPrinting()
	goFile.detectMapsAllocatedWithNew()
	goFile.detectEmptyIfBody()
	goFile.detectEmptyElseBody()
	goFile.detectEmptyForBody()
	goFile.detectGoToStatements()
	goFile.detectRaceInGoRoutine()
	goFile.detectIgnoredErrors()
	goFile.detectStaticCondition()
	goFile.detectStringMethodCallingItself()
	goFile.detectEarlyReturn()
	goFile.detectReturnBeforeElse()
	goFile.detectBufferNotFlushed()
}

type Violation struct {
	Type    Rule
	SrcLine int
	SrcPath string
}

func (violation *Violation) String() string {
	return fmt.Sprintf("%s on line %d", violation.Type, violation.SrcLine)
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

type file struct {
	goFile     *ast.File
	fileSet    *token.FileSet
	typeInfo   *types.Info
	violations []*Violation

	filePath string
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.goFile)
}

func (f *file) AddViolation(tokPosition token.Pos, violationType Rule) {
	srcLine := getSourceCodeLineNumber(f.fileSet, tokPosition)
	violation := &Violation{
		Type:    violationType,
		SrcLine: srcLine,
		SrcPath: f.filePath,
	}
	f.violations = append(f.violations, violation)
}

func (f *file) detectFmtPrinting() {
	f.walk(func(node ast.Node) bool {
		fmtMethods := []string{"Print", "Println", "Printf"}

		switch t := node.(type) {
		case *ast.SelectorExpr:
			if packName, ok := t.X.(*ast.Ident); ok {
				for _, method := range fmtMethods {
					if packName.Name == "fmt" && t.Sel.Name == method {
						f.AddViolation(t.Pos(), FMT_PRINTING)
						return false
					}
				}
			}
		}
		return true
	})
}

func (f *file) detectMapsAllocatedWithNew() {
	f.walk(func(node ast.Node) bool {
		newKeyword := "new"

		switch t := node.(type) {
		case *ast.CallExpr:
			if value, ok := t.Fun.(*ast.Ident); ok {
				if ok && value.Name == newKeyword {
					if _, ok := t.Args[0].(*ast.MapType); ok {
						f.AddViolation(t.Pos(), MAP_ALLOCATED_WITH_NEW)
						return false
					}
				}
			}
		}
		return true
	})
}

func (f *file) detectEmptyIfBody() {
	f.walk(func(node ast.Node) bool {
		if ifStmt, ok := node.(*ast.IfStmt); ok {
			if len(ifStmt.Body.List) == 0 {
				f.AddViolation(ifStmt.Pos(), EMPTY_IF_BODY)
				return false
			}
		}
		return true
	})
}

func (f *file) detectEmptyElseBody() {
	f.walk(func(node ast.Node) bool {
		if ifStmt, ok := node.(*ast.IfStmt); ok {
			if elseBody, ok := ifStmt.Else.(*ast.BlockStmt); ok {
				if len(elseBody.List) == 0 {
					f.AddViolation(elseBody.Pos(), EMPTY_ELSE_BODY)
					return false
				}
			}
		}
		return true
	})
}

func (f *file) detectEmptyForBody() {
	f.walk(func(node ast.Node) bool {
		if forStmt, ok := node.(*ast.ForStmt); ok {
			if len(forStmt.Body.List) == 0 {
				f.AddViolation(forStmt.Pos(), EMPTY_FOR_BODY)
				return false
			}
		}
		return true
	})
}

func (f *file) detectGoToStatements() {
	f.walk(func(node ast.Node) bool {
		if branchStmt, ok := node.(*ast.BranchStmt); ok {
			if branchStmt.Label != nil {
				f.AddViolation(branchStmt.Pos(), GOTO_USED)
				return false
			}
		}
		return true
	})
}

func (f *file) detectRaceInGoRoutine() {
	f.walk(func(node ast.Node) bool {
		if goStmt, ok := node.(*ast.GoStmt); ok {

			if goFunc, ok := goStmt.Call.Fun.(*ast.FuncLit); ok {
				goFuncParams := goFunc.Type.Params.List

				for _, goFuncBodyStmt := range goFunc.Body.List {
					if exprStmt, ok := goFuncBodyStmt.(*ast.ExprStmt); ok {

						if !validateParams(exprStmt, goFuncParams) {
							f.AddViolation(node.Pos(), RACE_CONDITION)
							return false
						}

					}
				}
			}
		}
		return true
	})
}

func (f *file) detectIgnoredErrors() {
	errorType := "error"
	ignored := false

	f.walk(func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.FuncDecl:
			ignored = ruleIgnored(ERROR_IGNORED, t.Doc)

		case *ast.AssignStmt:
			if !ignored {
				var ignoredReturnIndex []int
				for index, expr := range t.Lhs {
					if varName, ok := expr.(*ast.Ident); ok {
						if varName.Name == "_" {
							ignoredReturnIndex = append(ignoredReturnIndex, index)
						}
					}
				}
				if len(ignoredReturnIndex) != 0 {
					if tv, ok := f.typeInfo.Types[t.Rhs[0]]; ok {
						if tuple, ok := tv.Type.(*types.Tuple); ok {
							for _, returnIndex := range ignoredReturnIndex {
								if returnIndex <= tuple.Len() && tuple.At(returnIndex).Type().String() == errorType {
									f.AddViolation(t.Pos(), ERROR_IGNORED)
									return false
								}
							}
						}
					}
				}
			}
			return false

		case *ast.CallExpr:
			if !ignored {
				if tv, ok := f.typeInfo.Types[t]; ok {
					if name, ok := tv.Type.(*types.Named); ok {
						if name.String() == errorType {
							f.AddViolation(t.Pos(), ERROR_IGNORED)
						}
					} else if tuple, ok := tv.Type.(*types.Tuple); ok {
						for i := 0; i < tuple.Len(); i++ {
							if tuple.At(i).Type().String() == errorType {
								f.AddViolation(t.Pos(), ERROR_IGNORED)
							}
						}
					}
				}
			}
		}

		return true
	})
}

func (goFile *file) detectStringMethodCallingItself() {
	goFile.walk(func(node ast.Node) bool {
		return true
	})
}

func (goFile *file) detectStaticCondition() {
	goFile.walk(func(node ast.Node) bool {
		//TODO
		return true
	})
}

func (goFile *file) detectEarlyReturn() {
	goFile.walk(func(node ast.Node) bool {

		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			bodyLength := len(funcDecl.Body.List) - 1
			for index, bodyStmt := range funcDecl.Body.List {

				if _, ok := bodyStmt.(*ast.ReturnStmt); ok && bodyLength > index {
					goFile.AddViolation(bodyStmt.Pos(), EARLY_RETURN)
					return false
				}
			}
		}
		return true
	})
}

func (goFile *file) detectReturnBeforeElse() {
	goFile.walk(func(node ast.Node) bool {
		if ifStmt, ok := node.(*ast.IfStmt); ok {
			bodyLength := len(ifStmt.Body.List) - 1

			if ifStmt.Else != nil && bodyLength >= 0 {
				if retStmt, ok := ifStmt.Body.List[bodyLength].(*ast.ReturnStmt); ok {
					goFile.AddViolation(retStmt.Pos(), NO_ELSE_RETURN)
				}
			}
		}
		return true
	})
}

func (goFile *file) detectBufferNotFlushed() {
	goFile.walk(func(node ast.Node) bool {
		// TODO
		return true
	})
}

func ruleIgnored(ruleToIgnore Rule, commentGroup *ast.CommentGroup) bool {
	suppressRule := regexp.MustCompile(`@SuppressRule\("(?P<Rule>\S+)"\)`)
	if commentGroup != nil {
		for _, comment := range commentGroup.List {
			if result := suppressRule.FindStringSubmatch(comment.Text); len(result) > 0 {
				if result[1] == ruleToIgnore.String() {
					return true
				}
			}
		}
	}
	return false
}

// validateParams checks that all statements in function body is referencing
// local variable (including function argument), and not referencing outer scope.
func validateParams(exprStmt *ast.ExprStmt, List []*ast.Field) bool {
	if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
		// It is a call expression.
		for _, argument := range callExpr.Args {
			// check all function arguments
			if argumentName, ok := argument.(*ast.Ident); ok {
				if !nameInFieldList(argumentName, List) {
					return false
				}
			}
		}
	}
	return true
}

// nameInFieldList check ifs 'fieldList' (field/method/parameter) contains the name 'name'.
func nameInFieldList(name *ast.Ident, fieldList []*ast.Field) bool {
	for _, field := range fieldList {
		if field != nil {
			// Field might be anonymous.
			for _, fieldName := range field.Names {
				if fieldName.Name == name.Name {
					return true
				}
			}
		}
	}
	return false
}

// getSourceCodeLineNumber returns the line number in the parsed source code, according to the tokens position.
func getSourceCodeLineNumber(fileSet *token.FileSet, position token.Pos) int {
	return fileSet.File(position).Line(position)
}
