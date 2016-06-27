// Copyright (c) 2015-2016 The GoAnalysis Authors. All rights reserved.
// Use of this source code is governed by the MIT license found in the
// LICENSE.txt file.

package linter

import (
	"fmt"
	"github.com/chrisbbe/GoAnalysis/analysis/linter/astutil"
	"github.com/chrisbbe/GoAnalysis/analysis/linter/racer"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"runtime/debug"
	"errors"
)

const PRINT_STACK = true

type Rule int

//Rules.
const (
	RACE_CONDITION Rule = iota
	FMT_PRINTING
	STRING_DEFINES_ITSELF
	NEVER_ALLOCATED_MAP_WITH_NEW
	ERROR_IGNORED
	EMPTY_IF_BODY
	EMPTY_ELSE_BODY
	EMPTY_FOR_BODY
	GOTO_USED
	CONDITION_EVALUATED_STATICALLY
	NO_BUFFERED_FLUSHING
)

var ruleStrings = [...]string{
	RACE_CONDITION:                 "RACE_CONDITION",
	FMT_PRINTING:                   "FMT_PRINTING",
	STRING_DEFINES_ITSELF:          "STRING_DEFINES_ITSELF",
	NEVER_ALLOCATED_MAP_WITH_NEW:   "NEVER_ALLOCATE_MAP_WITH_NEW",
	ERROR_IGNORED:                  "ERROR_IGNORED",
	EMPTY_IF_BODY:                  "EMPTY_IF_BODY",
	EMPTY_ELSE_BODY:                "EMPTY_ELSE_BODY",
	EMPTY_FOR_BODY:                 "EMPTY_FOR_BODY",
	GOTO_USED:                      "GOTO_USED",
	CONDITION_EVALUATED_STATICALLY: "CONDITION_EVALUATED_STATICALLY",
	NO_BUFFERED_FLUSHING:           "NO_BUFFERED_FLUSHING",
}

func (rule Rule) String() string {
	return ruleStrings[rule]
}

type Violation struct {
	SrcPath string
	SrcLine int
	Type    Rule
}

// visitor maintains the state of the AST traversal and holds the result of mistakes detected.
type visitor struct {
	goFile          string
	fileSet         *token.FileSet
	currentFunction *ast.FuncDecl
	currentFile     *ast.File
	currentImports  []*ast.ImportSpec
	mistakes        []*Violation
}

func (v *visitor) AddViolation(tokenPosition token.Pos, MistakeType Rule) {
	srcLine := getLineNumberInSourceCode(v.fileSet, tokenPosition)
	v.mistakes = append(v.mistakes, &Violation{SrcPath: v.goFile, SrcLine: srcLine, Type: MistakeType})
}

func (v *visitor) getImportPackagePath(packageName string) (string, error) {
	for _, imp := range v.currentImports {
		path := strings.Trim(imp.Path.Value, "\"") // Trim away leading and ending " .
		if filepath.Base(path) == packageName {
			return path, nil
		}
	}
	return "", fmt.Errorf("Could not find import path for package %s.", packageName)
}

//TODO: Omg, rewrite, make easier and prettier!
func (v *visitor) errorIgnored(leftHandSide []ast.Expr, rightHandSide []ast.Expr) bool {

	for _, expr := range rightHandSide {
		if value, ok := expr.(*ast.CallExpr); ok {
			if value1, ok := value.Fun.(*ast.SelectorExpr); ok {
				if value2, ok := value1.X.(*ast.Ident); ok {
					if packagePath, err := v.getImportPackagePath(value2.Name); err == nil {
						if functionResult, err := astutil.FindFuncDecl(value1.Sel.Name, packagePath); err == nil {
							for index, result := range functionResult {
								if value3, ok := result.Type.(*ast.Ident); ok && value3.Name == "error" {

									//log.Printf("Index size: %d\n", index)

									// Is error.
									if len(leftHandSide) > index {
										if value4, ok := leftHandSide[index].(*ast.Ident); ok && value4.Name == "_" {
											return true
										}
									}

								}
							}
						}
					}
				}
			} else if value1, ok := value.Fun.(*ast.Ident); ok && value1.Obj != nil && value1.Obj.Kind == ast.Fun {
				// Handle calls to local functions.
				if funcDecl, ok := value1.Obj.Decl.(*ast.FuncDecl); ok {

					if funcDecl.Type.Results != nil {
						for index, resultField := range funcDecl.Type.Results.List {
							if field, ok := resultField.Type.(*ast.Ident); ok && field.Name == "error" {
								if ident, ok := leftHandSide[index].(*ast.Ident); ok && ident.Name == "_" {
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

//FindCommonMistakes is the function creating the AST and
//initiating walking of the tree, and returning final result of analysis.
func DetectViolations(goFile string) (mistakes []*Violation, walkError error) {
	srcFile, err := ioutil.ReadFile(goFile)
	if err != nil {
		log.Printf("Error here(3): %s\n", err)
	}
	// Create the AST by parsing src.
	fileSet := token.NewFileSet() // positions are relative to fileSet
	file, err := parser.ParseFile(fileSet, "", srcFile, 0)
	if err != nil {
		return nil, fmt.Errorf("Parse error in %s (%s), Skipping file!\n", goFile, err)
	}

	//Visit() method may panic(), which catches in this deferred closure.
	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
			if PRINT_STACK {
				//TODO: Remove
				log.Print("Her går det gæern!")
				debug.PrintStack()
			}
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

	visitor := visitor{goFile: goFile, fileSet: fileSet}
	ast.Walk(&visitor, file)
	return visitor.mistakes, walkError
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch t := node.(type) {

		case *ast.File:
			v.currentFile = t
			v.currentImports = t.Imports

		case *ast.BranchStmt:
			if t.Label != nil {
				v.AddViolation(t.Pos(), GOTO_USED)
			}

		case *ast.GoStmt:
			if racer.DetectRaceInGoRoutine(t) {
				v.AddViolation(t.Pos(), RACE_CONDITION)
			}

		case *ast.SelectorExpr:
			value, ok := t.X.(*ast.Ident)
			if ok && value.Name == "fmt" && (t.Sel.Name == "Print" || t.Sel.Name == "Println" || t.Sel.Name == "Printf") {
				v.AddViolation(t.Pos(), FMT_PRINTING)
			}

			if v.currentFunction != nil && v.currentFunction.Name.Name == "String" && ok && (t.Sel.Name == "Sprintf") {
				v.AddViolation(t.Pos(), STRING_DEFINES_ITSELF)
			}

		case *ast.FuncDecl:
			v.currentFunction = t // Update visitor with current function decl.

		case *ast.CallExpr:
			//TODO: Refactor out in separate function.
			value, ok := t.Fun.(*ast.Ident)

			//Handle map allocation with New()
			if ok && value.Name == "new" {
				if _, mapType := t.Args[0].(*ast.MapType); mapType {
					v.AddViolation(t.Pos(), NEVER_ALLOCATED_MAP_WITH_NEW)
				}
			}

		case *ast.IfStmt:
			// Checking if If-body is empty.
			if len(t.Body.List) == 0 {
				v.AddViolation(t.Pos(), EMPTY_IF_BODY)
			}

			// Checking if Else-body is empty.
			if elseBody, ok := t.Else.(*ast.BlockStmt); ok {
				if len(elseBody.List) == 0 {
					v.AddViolation(elseBody.Pos(), EMPTY_ELSE_BODY)
				}
			}

			// Check if cond is 'true' or 'false'.
			if stringCond, ok := t.Cond.(*ast.Ident); ok {
				if strings.EqualFold(stringCond.Name, "true") || strings.EqualFold(stringCond.Name, "false") {
					v.AddViolation(t.Pos(), CONDITION_EVALUATED_STATICALLY)
				}
			}

		case *ast.ForStmt:
			if len(t.Body.List) == 0 {
				v.AddViolation(t.Pos(), EMPTY_FOR_BODY)
			}

		case *ast.AssignStmt:
			if v.errorIgnored(t.Lhs, t.Rhs) {
				v.AddViolation(t.Pos(), ERROR_IGNORED)
			}
		}
	}
	return v
}

func getLineNumberInSourceCode(fileSet *token.FileSet, position token.Pos) int {
	return fileSet.File(position).Line(position)
}
