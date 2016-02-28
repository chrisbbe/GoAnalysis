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

// Package 'ccomplexity' provides functions to measure cyclomatic complexity
// of Go source code files, both cyclomatic complexity per function and file is
// supported.
package ccomplexity

import (
	"testing"
	"io/ioutil"
)

type functionCC struct {
	Complexity   int
	FunctionName string
}

func TestSimpleCyclomaticComplexityFileLevel(t *testing.T) {
	srcFile := "./testcode/_HelloWorld.go"

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Error finding file %s!\n", srcFile)
	}

	expectedCyclomaticComplexity := GetCyclomaticComplexityFileLevel(sourceFile)
	actualCyclomaticComplexity := 3
	if actualCyclomaticComplexity != expectedCyclomaticComplexity {
		t.Fatalf("Cyclomatic complexity in file %s should be %d, not %d!", srcFile, actualCyclomaticComplexity, expectedCyclomaticComplexity)
	}
}

func TestSimpleComplexityFunctionLevel(t *testing.T) {
	srcFile := "./testcode/_HelloWorld.go"
	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("Error finding file %s!\n", srcFile)
	}

	expectedCyclomaticComplexity := GetCyclomaticComplexityFunctionLevel(sourceFile)
	actualCyclomaticComplexity := []functionCC{functionCC{FunctionName:"main", Complexity:0}}
	if len(expectedCyclomaticComplexity) != len(actualCyclomaticComplexity) {
		t.Fatalf("Number of functions should be %d, but are %d!\n", len(actualCyclomaticComplexity),
			len(expectedCyclomaticComplexity))
	}

	for i, eCC := range expectedCyclomaticComplexity {
		if actualCyclomaticComplexity[i].FunctionName != eCC.FunctionName {
			t.Errorf("Function nr.%d should be named %s, and not %s!\n", i, actualCyclomaticComplexity[i].FunctionName, eCC.FunctionName)
		}
		if actualCyclomaticComplexity[i].Complexity != eCC.GetCyclomaticComplexity() {
			t.Errorf("Function '%s' should have cyclomatic complexity %d, and not %d!", eCC.FunctionName, actualCyclomaticComplexity[i].Complexity, eCC.GetCyclomaticComplexity())
		}
	}
}
