// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package linter

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"regexp"
	"strings"

	"bytes"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity"
	"io/ioutil"
	"log"
	"path/filepath"
)

const CC_LIMIT = 10 // Upper limit of cyclomatic complexity measures.
const ERROR_OUTPUT_FILE = "GoAnalyzerError.log"

var errorFileLogger *log.Logger

// Set logging to file!.
func init() {
	errorLogFile, err := os.OpenFile(ERROR_OUTPUT_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	errorFileLogger = log.New(errorLogFile, "", log.Lshortfile)
}

type Rule int

//Rules.
const (
	RACE_CONDITION Rule = iota
	FMT_PRINTING
	STRING_CALLS_ITSELF
	MAP_ALLOCATED_WITH_NEW
	ERROR_IGNORED
	EMPTY_IF_BODY
	EMPTY_ELSE_BODY
	EMPTY_FOR_BODY
	GOTO_USED
	CONDITION_EVALUATED_STATICALLY
	BUFFER_NOT_FLUSHED
	RETURN_KILLS_CODE
	CYCLOMATIC_COMPLEXITY
)

var ruleStrings = map[Rule]string{
	RACE_CONDITION:                 "RACE_CONDITION",
	FMT_PRINTING:                   "FMT_PRINTING",
	STRING_CALLS_ITSELF:            "STRING_CALLS_ITSELF",
	MAP_ALLOCATED_WITH_NEW:         "MAP_ALLOCATED_WITH_NEW",
	ERROR_IGNORED:                  "ERROR_IGNORED",
	EMPTY_IF_BODY:                  "EMPTY_IF_BODY",
	EMPTY_ELSE_BODY:                "EMPTY_ELSE_BODY",
	EMPTY_FOR_BODY:                 "EMPTY_FOR_BODY",
	GOTO_USED:                      "GOTO_USED",
	CONDITION_EVALUATED_STATICALLY: "CONDITION_EVALUATED_STATICALLY",
	BUFFER_NOT_FLUSHED:             "NO_BUFFERED_FLUSHING",
	RETURN_KILLS_CODE:              "RETURN_KILLS_CODE",
	CYCLOMATIC_COMPLEXITY:          "CYCLOMATIC_COMPLEXITY",
}

type GoFile struct {
	FilePath        string
	LinesOfCode     int
	LinesOfComments int
	Violations      []*Violation

	goFileNode *ast.File
	fileSet    *token.FileSet
	typeInfo   *types.Info

	typeErrorLogFile *os.File
}

type GoPackage struct {
	Violations []*GoFile

	Path string
	Pack *ast.Package

	fileSet  *token.FileSet
	typeInfo *types.Info
}

func (goPackage *GoPackage) GetFileNodes() (goFiles []*ast.File) {
	for _, goFile := range goPackage.Pack.Files {
		goFiles = append(goFiles, goFile)
	}
	return
}

type Violation struct {
	Type        Rule
	Description string
	SrcLine     int
}

func (rule Rule) String() string {
	return ruleStrings[rule]
}

func (violation *Violation) String() string {
	return fmt.Sprintf("%s (Line %d) - %s.", violation.Type, violation.SrcLine, violation.Description)
}

// Marshal Rule string value instead of int value.
func (rule Rule) MarshalText() ([]byte, error) {
	var buffer bytes.Buffer
	if _, err := buffer.WriteString(rule.String()); err != nil {
		return buffer.Bytes(), err
	}
	return buffer.Bytes(), nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func getAllDirectories(searchDir string) (map[string]map[string]*ast.Package, *token.FileSet) {
	dirPackages := make(map[string]map[string]*ast.Package)
	fileSet := token.NewFileSet()

	err := filepath.Walk(searchDir, func(folderPath string, f os.FileInfo, err error) error {
		if f.IsDir() {
			pks, firstErr := parser.ParseDir(fileSet, folderPath, nil, parser.ParseComments)
			if firstErr != nil {
				return firstErr
			}
			dirPackages[folderPath] = pks
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Parsing error: %s", err)
	}
	return dirPackages, fileSet
}

func DetectViolations(goSourceDir string) (goPackageViolations []*GoPackage, err error) {
	if isDir, err := isDirectory(goSourceDir); err == nil && !isDir {
		return goPackageViolations, fmt.Errorf("%s is not a directory", goSourceDir)
	} else if err != nil {
		return goPackageViolations, err
	}

	sourceDirPackages, fileSet := getAllDirectories(goSourceDir)

	for path, pkgs := range sourceDirPackages {
		for _, packNode := range pkgs {

			goPackage := &GoPackage{
				Path:    path,
				Pack:    packNode,
				fileSet: fileSet,
			}

			goPackage.measureCyclomaticComplexity(CC_LIMIT)
			goPackage.detectBugsAndCodeSmells()

			if len(goPackage.Violations) > 0 {
				goPackageViolations = append(goPackageViolations, goPackage)
			}
		}
	}

	return goPackageViolations, nil
}

// countLinesInFile counts number of lines which are code and comments and returns the result.
func (goFile *GoFile) countLinesInFile() {
	file, err := os.Open(goFile.FilePath)
	if err != nil {
		panic(err)
	}
	// Close file after reading!
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*/") {
				goFile.LinesOfComments++
			} else {
				goFile.LinesOfCode++
			}
		}
	}
}

func (goPackage *GoPackage) measureCyclomaticComplexity(upperLimit int) error {
	for filePath, goFile := range goPackage.Pack.Files {
		srcFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		complexity, err := ccomplexity.GetCyclomaticComplexityFunctionLevel(filePath, srcFile)
		if err != nil {
			return err
		}

		for _, funCC := range complexity {
			if funCC.Complexity > upperLimit {

				goFile := &GoFile{
					FilePath:   filePath,
					goFileNode: goFile,
				}
				goFile.countLinesInFile()

				goFile.Violations = append(goFile.Violations, &Violation{
					Type:    CYCLOMATIC_COMPLEXITY,
					SrcLine: funCC.SrcLine,
					Description: fmt.Sprintf("Cyclomatic complexity in %s() is %d, upper limit is %d.",
						funCC.Name, funCC.Complexity, upperLimit),
				})
				goPackage.Violations = append(goPackage.Violations, goFile)
			}
		}
	}
	return nil
}

func (goPackage *GoPackage) detectBugsAndCodeSmells() error {
	conf := types.Config{
		Importer:                 importer.Default(),
		FakeImportC:              true,
		DisableUnusedImportCheck: true,
	}

	goPackage.typeInfo = &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	if _, err := conf.Check(goPackage.Pack.Name, goPackage.fileSet, goPackage.GetFileNodes(), goPackage.typeInfo); err != nil {
		errorFileLogger.Printf("Error: %s", err)
	}

	goPackage.Analyze()

	return nil
}

func (goPackage *GoPackage) Analyze() {
	for filePath, file := range goPackage.Pack.Files {
		goFile := &GoFile{
			FilePath:   filePath,
			goFileNode: file,
			typeInfo:   goPackage.typeInfo,
			fileSet:    goPackage.fileSet,
		}

		goFile.Analyse()

		if len(goFile.Violations) > 0 {
			goFile.countLinesInFile()
			goPackage.Violations = append(goPackage.Violations, goFile)
		}
	}
}

// Analyse fires off all detection algorithms on the goFile.
func (goFile *GoFile) Analyse() {
	goFile.detectFmtPrinting()
	goFile.detectMapsAllocatedWithNew()
	goFile.detectEmptyIfBody()
	goFile.detectEmptyElseBody()
	goFile.detectEmptyForBody()
	goFile.detectGoToStatements()
	goFile.detectRaceInGoRoutine()
	goFile.detectIgnoredErrors()
	goFile.detectStaticCondition()
	goFile.detectRecursiveStringMethods()
	goFile.detectReturnKillingCode()
	goFile.detectBufferNotFlushed()
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

func (f *GoFile) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.goFileNode)
}

// ruleIgnored inspects the commentGroup after a @SuppressRule("RULE_NAME") annotation
// corresponding to the ruleToIgnore specified. Return TRUE if found, ELSE otherwise.
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

// AddViolation adds a new violation as specified trough its argument to the GoFile list of violations.
func (f *GoFile) AddViolation(tokPosition token.Pos, violationType Rule, description string) *Violation {
	srcLine := getSourceCodeLineNumber(f.fileSet, tokPosition)
	violation := &Violation{
		Type:        violationType,
		Description: description,
		SrcLine:     srcLine,
	}
	f.Violations = append(f.Violations, violation)
	return violation
}

// Detect violations of rule: FMT_PRINTING.
func (goFile *GoFile) detectFmtPrinting() {
	ignored := false
	goFile.walk(func(node ast.Node) bool {
		fmtMethods := []string{"Print", "Println", "Printf"}

		switch t := node.(type) {
		case *ast.FuncDecl:
			ignored = ruleIgnored(FMT_PRINTING, t.Doc)
		case *ast.SelectorExpr:
			if !ignored {
				if packName, ok := t.X.(*ast.Ident); ok {
					for _, method := range fmtMethods {
						if packName.Name == "fmt" && t.Sel.Name == method {
							goFile.AddViolation(
								t.Pos(),
								FMT_PRINTING,
								fmt.Sprint("Printing from the fmt package are not synchronized and usually intended for"+
									" debugging purposes. Consider to use the log package!"),
							)
							return false
						}
					}
				}
			}
		}
		return true
	})
}

// Detect violations of rule: MAP_ALLOCATED_WITH_NEW.
func (goFile *GoFile) detectMapsAllocatedWithNew() {
	goFile.walk(func(node ast.Node) bool {
		newKeyword := "new"

		switch t := node.(type) {
		case *ast.CallExpr:
			if value, ok := t.Fun.(*ast.Ident); ok {
				if ok && value.Name == newKeyword {
					if _, ok := t.Args[0].(*ast.MapType); ok {
						goFile.AddViolation(
							t.Pos(),
							MAP_ALLOCATED_WITH_NEW,
							fmt.Sprint("Maps must be initialized with make(), new() allocates a nil map causing runtime "+
								"panic on write operations!"),
						)
						return false
					}
				}
			}
		}
		return true
	})
}

// Detect violations of rule: EMPTY_IF_BODY.
func (goFile *GoFile) detectEmptyIfBody() {
	goFile.walk(func(node ast.Node) bool {
		if ifStmt, ok := node.(*ast.IfStmt); ok {
			if len(ifStmt.Body.List) == 0 {
				goFile.AddViolation(
					ifStmt.Pos(),
					EMPTY_IF_BODY,
					fmt.Sprint("If body is empty, wasteful to not do anything with the if condition."),
				)
				return false
			}
		}
		return true
	})
}

// Detect violations of rule: EMPTY_ELSE_BODY.
func (goFile *GoFile) detectEmptyElseBody() {
	goFile.walk(func(node ast.Node) bool {
		if ifStmt, ok := node.(*ast.IfStmt); ok {
			if elseBody, ok := ifStmt.Else.(*ast.BlockStmt); ok {
				if len(elseBody.List) == 0 {
					goFile.AddViolation(
						elseBody.Pos(),
						EMPTY_ELSE_BODY,
						fmt.Sprint("ELse body is empty, wasteful to not do anything with the else condition."),
					)
					return false
				}
			}
		}
		return true
	})
}

// Detect violations of rule: EMPTY_FOR_BODY.
func (goFile *GoFile) detectEmptyForBody() {
	goFile.walk(func(node ast.Node) bool {
		if forStmt, ok := node.(*ast.ForStmt); ok {
			if len(forStmt.Body.List) == 0 {
				goFile.AddViolation(
					forStmt.Pos(),
					EMPTY_FOR_BODY,
					fmt.Sprint("For body is empty, wasteful to not do anything with the for condition."),
				)
				return false
			}
		}
		return true
	})
}

// Detect violations of rule: GOTO_USED.
func (goFile *GoFile) detectGoToStatements() {
	goFile.walk(func(node ast.Node) bool {
		if branchStmt, ok := node.(*ast.BranchStmt); ok {
			if branchStmt.Label != nil {
				goFile.AddViolation(
					branchStmt.Pos(),
					GOTO_USED,
					fmt.Sprint("Please dont use GOTO statements, they lead to spagehetti code!"),
				)
				return false
			}
		}
		return true
	})
}

// Detect violations of rule: RACE_CONDITION.
func (goFile *GoFile) detectRaceInGoRoutine() {
	goFile.walk(func(node ast.Node) bool {
		if goStmt, ok := node.(*ast.GoStmt); ok {

			if goFunc, ok := goStmt.Call.Fun.(*ast.FuncLit); ok {
				goFuncParams := goFunc.Type.Params.List

				for _, goFuncBodyStmt := range goFunc.Body.List {
					if exprStmt, ok := goFuncBodyStmt.(*ast.ExprStmt); ok {

						if !validateParams(exprStmt, goFuncParams) {
							goFile.AddViolation(
								node.Pos(),
								RACE_CONDITION,
								fmt.Sprint("Loop iterator variables must be passed as argument to Goroutine, not referenced."),
							)
							return false
						}

					}
				}
			}
		}
		return true
	})
}

// Detect violations of rule: RETURN_KILLS_CODE.
func (goFile *GoFile) detectReturnKillingCode() {
	goFile.walk(func(node ast.Node) bool {

		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			if funcDecl.Body != nil {
				bodyLength := len(funcDecl.Body.List) - 1
				for index, bodyStmt := range funcDecl.Body.List {

					if _, ok := bodyStmt.(*ast.ReturnStmt); ok && bodyLength > index {
						goFile.AddViolation(
							bodyStmt.Pos(),
							RETURN_KILLS_CODE,
							fmt.Sprint("Code is dead because of return! There is no possible execution path to the code below in "+
								"this scope!"),
						)
						return false
					}
				}
			}
		}
		return true
	})
}

// Detect violations of rule: ERROR_IGNORE.
func (goFile *GoFile) detectIgnoredErrors() {
	errorType := "error"
	ignored := false
	var returnResults []ast.Expr
	var rightHandSideCallExpr []*ast.CallExpr // Holds CallExpr taken part in AssignStmt, avoid checking these CallExpr.

	goFile.walk(func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.FuncDecl:
			ignored = ruleIgnored(ERROR_IGNORED, t.Doc)

		case *ast.ReturnStmt:
			// Hold the list of return result expressions.
			returnResults = t.Results

		case *ast.AssignStmt:
			// Save CalLExpr that is part of Rhs.
			if callExpr, ok := t.Rhs[0].(*ast.CallExpr); ok {
				rightHandSideCallExpr = append(rightHandSideCallExpr, callExpr)
			}

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
					if tv, ok := goFile.typeInfo.Types[t.Rhs[0]]; ok {
						if tuple, ok := tv.Type.(*types.Tuple); ok {
							for _, returnIndex := range ignoredReturnIndex {
								if returnIndex <= tuple.Len() && tuple.At(returnIndex).Type().String() == errorType {
									goFile.AddViolation(
										t.Pos(),
										ERROR_IGNORED,
										fmt.Sprint("Never ignore erros, ignoring them can lead to program crashes"),
									)
									return false
								}
							}
						}
					}
				}
			}
			return true

		case *ast.CallExpr:
			if !ignored {
				// Loop through return result expression to check whether CallExpr is part of return.
				// Flagging errors not assigned to variable in return statement is wrong!
				for _, returnResult := range returnResults {
					if returnCallExpr, ok := returnResult.(*ast.CallExpr); ok {
						if returnCallExpr == t {
							// Call Expression is in return, stop looking after calls on functions returning error.
							return false
						}
					}
				}

				// We don't want to check CallExpressions that is part of AssignStmt.
				for _, rhsCallExpr := range rightHandSideCallExpr {
					if t == rhsCallExpr {
						return false
					}
				}

				// CallExpr is not part of return result, check further!
				if tv, ok := goFile.typeInfo.Types[t]; ok {
					if name, ok := tv.Type.(*types.Named); ok {

						if name.String() == errorType {
							goFile.AddViolation(
								t.Pos(),
								ERROR_IGNORED,
								fmt.Sprint("Never ignore erros, ignoring them can lead to program crashes"),
							)
							return false
						}
					} else if tuple, ok := tv.Type.(*types.Tuple); ok {
						for i := 0; i < tuple.Len(); i++ {
							if tuple.At(i).Type().String() == errorType {
								goFile.AddViolation(
									t.Pos(),
									ERROR_IGNORED,
									fmt.Sprint("Never ignore erros, ignoring them can lead to program crashes"),
								)
								return false
							}
						}
					}
				}
			}
		}
		return true
	})
}

func (goFile *GoFile) detectStaticCondition() {
	var booleanOperand = map[token.Token]bool{
		token.LAND: true, // &&
		token.LOR:  true, // ||
		token.LEQ:  true, // <=
		token.GEQ:  true, // >=
		token.NOT:  true, // !
		token.GTR:  true, // >
		token.LSS:  true, // <
		token.EQL:  true, // ==
		token.NEQ:  true, // !=
	}
	var lastViolation *Violation

	goFile.walk(func(node ast.Node) bool {
		switch t := node.(type) {

		case *ast.IfStmt:
			if cond, ok := t.Cond.(*ast.Ident); ok && (cond.Name == "true" || cond.Name == "false") {
				// Catch obvious static conditions.
				lastViolation = goFile.AddViolation(
					t.Pos(),
					CONDITION_EVALUATED_STATICALLY,
					fmt.Sprintf("Condition will always be %s", cond),
				)
			}

		case *ast.BinaryExpr:
			if booleanOperand[t.Op] {
				if _, ok := t.X.(*ast.BasicLit); ok {
					if _, ok := t.Y.(*ast.BasicLit); ok {
						// We don't want to add flag the same violation more than once in composed expressions eg. (x || y && k)
						// where x,y and k is basic literals.
						if lastViolation != nil && lastViolation.SrcLine != getSourceCodeLineNumber(goFile.fileSet, t.Pos()) {
							lastViolation = goFile.AddViolation(
								t.Pos(),
								CONDITION_EVALUATED_STATICALLY,
								fmt.Sprint("Both left and right operand is basic literal that can be eveualted at compile time"),
							)
						}
					}
				}
			}
			if t.Op == token.EQL || t.Op == token.NEQ {
				// If this is a comparison against nil, find the other operand.
				var other ast.Expr
				if goFile.typeInfo.Types[t.X].IsNil() {
					other = t.Y
				} else if goFile.typeInfo.Types[t.Y].IsNil() {
					other = t.X
				}

				// Find the object.
				var obj types.Object

				if other != nil {
					switch v := other.(type) {
					case *ast.Ident:
						obj = goFile.typeInfo.Uses[v]
					case *ast.SelectorExpr:
						obj = goFile.typeInfo.Uses[v.Sel]
					}
				}

				if obj != nil {
					if _, ok := obj.(*types.Func); ok {
						lastViolation = goFile.AddViolation(
							t.Pos(),
							CONDITION_EVALUATED_STATICALLY,
							fmt.Sprintf("Comparison of function %s is always %v", obj.Name(), t.Op == token.NEQ),
						)
					}
				}
			}

		}
		return true
	})
}

// TODO: Implement.
func (goFile *GoFile) detectRecursiveStringMethods() {
	goFile.walk(func(node ast.Node) bool {
		return false
	})
}

//TODO : Is it possible to actually do this without escape analysis?
func (goFile *GoFile) detectBufferNotFlushed() {
	goFile.walk(func(node ast.Node) bool {
		return true
	})
}
