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
	"io/ioutil"
	"testing"
	"fmt"
)

func verifyCyclomaticComplexity(expectedComplexity []*FunctionComplexity, correctComplexity []FunctionComplexity) error {
	if len(expectedComplexity) != len(correctComplexity) {
		return fmt.Errorf("Number of functions should be %d, but are %d!\n", len(expectedComplexity), len(correctComplexity))
	}

	for index, expectedComplexity := range expectedComplexity {
		if expectedComplexity.Name != correctComplexity[index].Name {
			return fmt.Errorf("Function nr. %d name should be %s, ant not %s.\n", index, correctComplexity[index].Name,
				expectedComplexity.Name)
		}

		if expectedComplexity.Complexity != correctComplexity[index].Complexity {
			fmt.Errorf("Function nr. %d (%s) should have cyclomatic complexity %d, but has %d!\n", index,
				expectedComplexity.Name, correctComplexity[index].Complexity, expectedComplexity.Complexity)
		}
	}
	return nil
}

func TestSimpleComplexityFunctionLevel(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_HelloWorld.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := GetCyclomaticComplexityFunctionLevel(srcFile)
	if err != nil {
		t.Fatal(err)
	}
	correctCyclomaticComplexity := []FunctionComplexity{
		FunctionComplexity{Name: "main", Complexity: 1},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}

func TestMethodComplexity(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_swap.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := GetCyclomaticComplexityFunctionLevel(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	correctCyclomaticComplexity := []FunctionComplexity{
		FunctionComplexity{Name: "main", Complexity: 1},
		FunctionComplexity{Name: "swap", Complexity: 1},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}

func TestSwitcherComplexity(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_switcher.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := GetCyclomaticComplexityFunctionLevel(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	correctCyclomaticComplexity := []FunctionComplexity{
		FunctionComplexity{Name: "main", Complexity: 1},
		FunctionComplexity{Name: "monthNumberToString", Complexity: 15},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}

func TestGreatestCommonDivisorComplexity(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_gcd.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := GetCyclomaticComplexityFunctionLevel(srcFile)
	if err != nil {
		t.Fatal(err)
	}
	correctCyclomaticComplexity := []FunctionComplexity{
		FunctionComplexity{Name: "gcd", Complexity: 2},
		FunctionComplexity{Name: "main", Complexity: 1},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}
