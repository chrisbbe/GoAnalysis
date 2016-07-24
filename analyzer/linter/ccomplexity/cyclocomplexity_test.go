// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package ccomplexity

import (
	"fmt"
	"io/ioutil"
	"testing"
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
			return fmt.Errorf("Function nr. %d (%s) should have cyclomatic complexity %d, but has %d!\n", index,
				expectedComplexity.Name, correctComplexity[index].Complexity, expectedComplexity.Complexity)
		}
	}
	return nil
}

func TestSimpleComplexityFunctionLevel(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_helloworld.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := GetCyclomaticComplexityFunctionLevel(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	correctCyclomaticComplexity := []FunctionComplexity{
		{Name: "main", Complexity: 1},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}

func TestIfElseComplexityFunctionLevel(t *testing.T) {
	srcFile, err := ioutil.ReadFile("./testcode/_ifelse.go")
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := GetCyclomaticComplexityFunctionLevel(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedCyclomaticComplexity[0].ControlFlowGraph.Draw("ifelse")

	correctCyclomaticComplexity := []FunctionComplexity{
		{Name: "main", Complexity: 2},
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
		{Name: "main", Complexity: 1},
		{Name: "swap", Complexity: 1},
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
		{Name: "main", Complexity: 1},
		{Name: "monthNumberToString", Complexity: 14},
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
		{Name: "gcd", Complexity: 2},
		{Name: "main", Complexity: 1},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}
