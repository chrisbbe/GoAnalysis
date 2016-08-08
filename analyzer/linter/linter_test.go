// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package linter_test

import (
	"fmt"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter"
	"testing"
)

// actualViolation is a type used in testing, to represent correct violation
// that should be detected.
type actualViolation struct {
	SrcLine int
	Type    linter.Rule
}

// verifyViolations checks the list of expected Violations with the list of actual Violations.
// Errors found is reported and logged through t.
func verifyViolations(expectedViolations []*linter.Violation, actualViolations []actualViolation) error {
	if len(expectedViolations) != len(actualViolations) {
		return fmt.Errorf("Number of violations should be %d, but are %d!\n", len(actualViolations), len(expectedViolations))
	}

	for index, expectedMistake := range expectedViolations {
		if actualViolations[index].Type != expectedMistake.Type {
			return fmt.Errorf("Violation should be of type %s, and not type %s!\n", actualViolations[index].Type, expectedMistake.Type)
		}
		if actualViolations[index].SrcLine != expectedMistake.SrcLine {
			return fmt.Errorf("Violation (%s) should be found on line %d, and not on line %d!\n", actualViolations[index].Type, actualViolations[index].SrcLine, expectedMistake.SrcLine)
		}
	}
	return nil
}

// Printing from fmt package is not thread safe and should be avoided in production and detected!
// Testing rule: FMT_PRINTING
func TestDetectionOfPrintingFromFmtPackage(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/fmtprinting")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 11, Type: linter.FMT_PRINTING},
		{SrcLine: 12, Type: linter.FMT_PRINTING},
		{SrcLine: 13, Type: linter.FMT_PRINTING},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: NEVER_ALLOCATED_MAP_WITH_NEW
// Allocating maps with new returns a nil pointer, therefor one should use make.
func TestDetectionOfAllocatingMapWithNew(t *testing.T) {
	expectedViolations, _ := linter.DetectViolations("./testcode/newmap")
	actualViolations := []actualViolation{
		{SrcLine: 11, Type: linter.MAP_ALLOCATED_WITH_NEW},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}

}

// Testing rule: RACE_CONDITION
// Races will occur, since multiple Go-routines will share the same counter variable.
func TestDetectionOfRacesInLoopClosures(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/threadlooping")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 18, Type: linter.RACE_CONDITION},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: EMPTY_IF_BODY
// Empty if-bodies are unnecessary and ineffective.
func TestDetectionOfEmptyIfBody(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/emptyifbody")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 11, Type: linter.EMPTY_IF_BODY},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: EMPTY_ELSE_BODY
// Empty else-bodies are unnecessary and ineffective.
func TestDetectionOfEmptyElseBody(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/emptyelsebody")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 16, Type: linter.EMPTY_ELSE_BODY},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: EMPTY_FOR_BODY
// Empty for-bodies are unnecessary and ineffective.
func TestDetectionOfEmptyForBody(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/emptyforbody")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 8, Type: linter.EMPTY_FOR_BODY},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: RETURN_KILLS_CODE
// One should never return unconditionally, except from the last statement in a func or method.
func TestDetectionOfEarlyReturn(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/earlyreturn")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 13, Type: linter.RETURN_KILLS_CODE},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: GOTO_USED
// Jumping around in the code using GOTO (including BREAK, CONTINUE, GOTO, FALLTHROUGH)
// is considered confusing and harmful.
func TestDetectionOfGoTo(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/goto")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 12, Type: linter.GOTO_USED},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: GOTO_USED
// Jumping around in the code using GOTO (including BREAK, CONTINUE, GOTO, FALLTHROUGH)
// is considered confusing and harmful.
func TestDetectionOfLabeledBranching(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/labeledbranch")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 15, Type: linter.GOTO_USED},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: ERROR_IGNORED
// Errors should never be ignored, might lead to program crashes.
func TestDetectionOfIgnoredErrors(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/errorignored")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 16, Type: linter.ERROR_IGNORED},
		{SrcLine: 17, Type: linter.ERROR_IGNORED},
		{SrcLine: 22, Type: linter.ERROR_IGNORED},
		{SrcLine: 27, Type: linter.ERROR_IGNORED},
		{SrcLine: 35, Type: linter.ERROR_IGNORED},
		{SrcLine: 51, Type: linter.ERROR_IGNORED},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: CONDITION_EVALUATED_STATICALLY
// Condition that can be evaluated statically are wasted and performance-reducing.
func TestDetectionOfConditionEvaluatedStatically(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/staticconditions")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 14, Type: linter.CONDITION_EVALUATED_STATICALLY},
		{SrcLine: 18, Type: linter.CONDITION_EVALUATED_STATICALLY},
		{SrcLine: 22, Type: linter.CONDITION_EVALUATED_STATICALLY},
		{SrcLine: 26, Type: linter.CONDITION_EVALUATED_STATICALLY},
		//{SrcLine: 32, Type: linter.CONDITION_EVALUATED_STATICALLY}, TODO: Possible to statically trace a variable value?
		{SrcLine: 38, Type: linter.CONDITION_EVALUATED_STATICALLY},
		{SrcLine: 41, Type: linter.CONDITION_EVALUATED_STATICALLY},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// GitHub Issue #4 list a scenario where the tool detects a false positive of rule ERROR_IGNORED.
// This test verifies correction of the behaviour.
func TestDetectionOfErrorIgnoredInReturn(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/errorinreturn")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{}

	if len(expectedViolations) > 0 {
		// Only verify violations if there is more than one file containing violations!
		if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDetectionOfHighCyclomatiComplexity(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/cyclomaticomplexity")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 13, Type: linter.CYCLOMATIC_COMPLEXITY},
	}

	if len(expectedViolations) > 0 {
		// Only verify violations if there is more than one file containing violations!
		if err := verifyViolations(expectedViolations[0].Violations[0].Violations, actualViolations); err != nil {
			t.Fatal(err)
		}
	}
}

//TODO: Implement!
/*
// Testing rule: STRING_METHOD_DEFINES_ITSELF
func TestDetectionOfStringMethodDefiningItself(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/stringmethod")
	if err != nil {
		t.Fatal(err)
	}

	correctViolations := []actualViolation{
		{SrcLine: 48, Type: linter.STRING_CALLS_ITSELF},
		{SrcLine: 58, Type: linter.STRING_CALLS_ITSELF},
		{SrcLine: 63, Type: linter.STRING_CALLS_ITSELF},
		{SrcLine: 69, Type: linter.STRING_CALLS_ITSELF},
	}

	if len(expectedViolations) <= 0 {
		t.Fatal("There is no functions containing violations.")
	}
	if err := verifyViolations(expectedViolations[0].Violations[0].Violations, correctViolations); err != nil {
		t.Fatal(err)
	}
}
*/
