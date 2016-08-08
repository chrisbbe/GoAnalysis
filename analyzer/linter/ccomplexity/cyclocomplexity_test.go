// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package ccomplexity_test

import (
	"fmt"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity"
	"io/ioutil"
	"testing"
)

func verifyCyclomaticComplexity(expectedComplexity []*ccomplexity.FunctionComplexity, correctComplexity []ccomplexity.FunctionComplexity) error {
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
	filePath := "./testcode/_helloworld.go"
	srcFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := ccomplexity.GetCyclomaticComplexityFunctionLevel(filePath, srcFile)
	if err != nil {
		t.Fatal(err)
	}

	correctCyclomaticComplexity := []ccomplexity.FunctionComplexity{
		{Name: "main", Complexity: 1},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}

func TestIfElseComplexityFunctionLevel(t *testing.T) {
	filePath := "./testcode/_ifelse.go"
	srcFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := ccomplexity.GetCyclomaticComplexityFunctionLevel(filePath, srcFile)
	if err != nil {
		t.Fatal(err)
	}

	correctCyclomaticComplexity := []ccomplexity.FunctionComplexity{
		{Name: "main", Complexity: 2},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}

func TestMethodComplexity(t *testing.T) {
	filePath := "./testcode/_swap.go"
	srcFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := ccomplexity.GetCyclomaticComplexityFunctionLevel(filePath, srcFile)
	if err != nil {
		t.Fatal(err)
	}

	correctCyclomaticComplexity := []ccomplexity.FunctionComplexity{
		{Name: "main", Complexity: 1},
		{Name: "swap", Complexity: 1},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}

func TestSwitcherComplexity(t *testing.T) {
	filePath := "./testcode/_switcher.go"
	srcFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := ccomplexity.GetCyclomaticComplexityFunctionLevel(filePath, srcFile)
	if err != nil {
		t.Fatal(err)
	}

	correctCyclomaticComplexity := []ccomplexity.FunctionComplexity{
		{Name: "main", Complexity: 1},
		{Name: "monthNumberToString", Complexity: 14},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}

func TestGreatestCommonDivisorComplexity(t *testing.T) {
	filePath := "./testcode/_gcd.go"
	srcFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	expectedCyclomaticComplexity, err := ccomplexity.GetCyclomaticComplexityFunctionLevel(filePath, srcFile)
	if err != nil {
		t.Fatal(err)
	}
	correctCyclomaticComplexity := []ccomplexity.FunctionComplexity{
		{Name: "gcd", Complexity: 2},
		{Name: "main", Complexity: 1},
	}

	if err := verifyCyclomaticComplexity(expectedCyclomaticComplexity, correctCyclomaticComplexity); err != nil {
		t.Error(err)
	}
}
