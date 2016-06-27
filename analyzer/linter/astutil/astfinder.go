package astutil

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const (
	STANDARD_PACKAGE_PATH_ENV = "GOROOT"
	CUSTOM_PACKAGE_PATH_ENV = "GOPATH"
	SOURCE_FOLDER = "/src/"
)

type VisitorFunc func(n ast.Node) ast.Visitor

func (f VisitorFunc) Visit(n ast.Node) ast.Visitor {
	return f(n)
}

type visitor struct {
	FunctionName   string       //Name of function we are looking up.
	FunctionResult []*ast.Field // Returning results from function.
}

// FunctionFound returns True if the visitor contains a
// result, indicating that the function and its return
// values is found. Returning False otherwise.
func (v *visitor) FunctionFound() bool {
	return v.FunctionResult != nil
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if v.FunctionFound() {
		// No need to search when result is found.
		return nil
	}

	if node != nil {
		switch t := node.(type) {

		case *ast.FuncDecl:
			if t.Name.Name == v.FunctionName {
				// Function found, extract function return result.
				if t.Type.Results != nil {
					v.FunctionResult = t.Type.Results.List
				}
			}
		}
	}
	return v
}

// FindFunction walks the AST tree for the program, looking for the function signature specified.
// Tries first to look for the function in the standard library before looking into the custom library.
// Requires both $GOROOT and $GOPATH to be set, otherwise an error will be returned.
func FindFuncDecl(functionName string, packagePath string) ([]*ast.Field, error) {
	visitor := &visitor{FunctionName: functionName}

	if env, ok := os.LookupEnv(STANDARD_PACKAGE_PATH_ENV); ok {
		//Look into the standard library.
		findFuncDecl(visitor, env + SOURCE_FOLDER, packagePath)
	} else {
		return nil, fmt.Errorf("Environment variable $%s not set!", STANDARD_PACKAGE_PATH_ENV)
	}

	if !visitor.FunctionFound() {
		// Try to look into custom packages.
		if env, ok := os.LookupEnv(CUSTOM_PACKAGE_PATH_ENV); ok {
			findFuncDecl(visitor, env + SOURCE_FOLDER, packagePath)
		} else {
			return nil, fmt.Errorf("Environment variable $%s not set!", CUSTOM_PACKAGE_PATH_ENV)
		}
	}
	return visitor.FunctionResult, nil
}

func findFuncDecl(v *visitor, packageRoot string, packagePath string) {
	packages, err := parser.ParseDir(token.NewFileSet(), packageRoot + packagePath, nil, 0)
	if err != nil {
		return //Dir not present.
	}

	lastPackage := strings.Split(packagePath, "/")[strings.Count(packagePath, "/")] //Extract package from path.
	if pack, ok := packages[lastPackage]; ok {
		ast.Walk(v, pack)
	}
}
