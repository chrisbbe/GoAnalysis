// Copyright (c) 2015-2016 The GoAnalysis Authors. All rights reserved.
// Use of this source code is governed by the MIT license found in the
// LICENSE.txt file.

// Package racer implements some race condition detection functions.
package racer

import "go/ast"

//TODO: MA-22, not agreeing with myself about this!
// DetectRaceInGoRoutine detect if params used in Go routine is passed properly to the the function,
// and not referenced to outer scope.
func DetectRaceInGoRoutine(goStmt *ast.GoStmt) bool {
	if goFunc, ok := goStmt.Call.Fun.(*ast.FuncLit); ok {
		goFuncParams := goFunc.Type.Params.List

		for _, goFuncBodyStmt := range goFunc.Body.List {
			if exprStmt, ok := goFuncBodyStmt.(*ast.ExprStmt); ok {

				if !validateParams(exprStmt, goFuncParams) {
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
