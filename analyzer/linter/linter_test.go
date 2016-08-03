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
	expectedViolations, err := linter.DetectViolations("./testcode/_fmtprinting.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 10, Type: linter.FMT_PRINTING},
		{SrcLine: 11, Type: linter.FMT_PRINTING},
		{SrcLine: 12, Type: linter.FMT_PRINTING},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: NEVER_ALLOCATED_MAP_WITH_NEW
// Allocating maps with new returns a nil pointer, therefor one should use make.
func TestDetectionOfAllocatingMapWithNew(t *testing.T) {
	expectedViolations, _ := linter.DetectViolations("./testcode/_newmap.go")
	actualViolations := []actualViolation{
		{SrcLine: 11, Type: linter.MAP_ALLOCATED_WITH_NEW},
	}

	if len(expectedViolations) != len(actualViolations) {
		t.Fatalf("Number of Violations should be %d, not %d!\n", len(actualViolations), len(expectedViolations))
	}

	for index, expectedMistake := range expectedViolations[0].Violations {
		if actualViolations[index].Type != expectedMistake.Type {
			t.Errorf("Error should be of type %s, and not type %s!\n", actualViolations[index].Type, expectedMistake.Type)
		}
		if actualViolations[index].SrcLine != expectedMistake.SrcLine {
			t.Errorf("Error should be found on line %d, and not on line %d!\n", actualViolations[index].SrcLine, expectedMistake.SrcLine)
		}
	}

}

// Testing rule: RACE_CONDITION
// Races will occur, since multiple Go-routines will share the same counter variable.
func TestDetectionOfRacesInLoopClosures(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_threadlooping.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 18, Type: linter.RACE_CONDITION},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: EMPTY_IF_BODY
// Empty if-bodies are unnecessary and ineffective.
func TestDetectionOfEmptyIfBody(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_emptyifbody.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 11, Type: linter.EMPTY_IF_BODY},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: EMPTY_ELSE_BODY
// Empty else-bodies are unnecessary and ineffective.
func TestDetectionOfEmptyElseBody(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_emptyelsebody.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 16, Type: linter.EMPTY_ELSE_BODY},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: EMPTY_FOR_BODY
// Empty for-bodies are unnecessary and ineffective.
func TestDetectionOfEmptyForBody(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_emptyforbody.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 8, Type: linter.EMPTY_FOR_BODY},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: RETURN_KILLS_CODE
// One should never return unconditionally, except from the last statement in a func or method.
func TestDetectionOfEarlyReturn(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_earlyreturn.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 13, Type: linter.RETURN_KILLS_CODE},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: NO_ELSE_RETURN
// When an If block contains a return statement, the Else block becoms unnecessary.
func TestDetectionOfReturnBeforeElse(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_elsereturn.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 14, Type: linter.RETURN_KILLS_CODE},
		{SrcLine: 22, Type: linter.RETURN_KILLS_CODE},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: GOTO_USED
// Jumping around in the code using GOTO (including BREAK, CONTINUE, GOTO, FALLTHROUGH)
// is considered confusing and harmful.
func TestDetectionOfGoTo(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_goto.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 12, Type: linter.GOTO_USED},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: GOTO_USED
// Jumping around in the code using GOTO (including BREAK, CONTINUE, GOTO, FALLTHROUGH)
// is considered confusing and harmful.
func TestDetectionOfLabeledBranching(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_labeledbranch.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 15, Type: linter.GOTO_USED},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: ERROR_IGNORED
// Errors should never be ignored, might lead to program crashes.
func TestDetectionOfIgnoredErrors(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_ignoreerror.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		{SrcLine: 16, Type: linter.ERROR_IGNORED},
		{SrcLine: 17, Type: linter.ERROR_IGNORED},
		{SrcLine: 22, Type: linter.ERROR_IGNORED},
		{SrcLine: 27, Type: linter.ERROR_IGNORED},
		{SrcLine: 30, Type: linter.ERROR_IGNORED},
		{SrcLine: 46, Type: linter.ERROR_IGNORED},
	}

	if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

/*
// Testing rule: STRING_METHOD_DEFINES_ITSELF
func TestDetectionOfStringMethodDefiningItself(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_stringmethod.go")
	if err != nil {
		t.Fatal(err)
	}

	for _, vio := range expectedViolations {
		log.Printf("- %s\n", vio)
	}

	correctViolations := []actualViolation{
		actualViolation{SrcLine: 48, Type: linter.STRING_DEFINES_ITSELF},
		actualViolation{SrcLine: 58, Type: linter.STRING_DEFINES_ITSELF},
		actualViolation{SrcLine: 63, Type: linter.STRING_DEFINES_ITSELF},
		actualViolation{SrcLine: 69, Type: linter.STRING_DEFINES_ITSELF},
	}

	if err := verifyViolations(expectedViolations, correctViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: CONDITION_EVALUATED_STATICALLY
// Condition that can be evaluated statically are wasted and performance-reducing.
func TestDetectionOfConditionEvaluatedStatically(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_staticconditions.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		actualViolation{SrcLine: 14, Type: linter.CONDITION_EVALUATED_STATICALLY},
		actualViolation{SrcLine: 18, Type: linter.CONDITION_EVALUATED_STATICALLY},
		actualViolation{SrcLine: 22, Type: linter.CONDITION_EVALUATED_STATICALLY},
		actualViolation{SrcLine: 26, Type: linter.CONDITION_EVALUATED_STATICALLY},
		actualViolation{SrcLine: 33, Type: linter.CONDITION_EVALUATED_STATICALLY},
		actualViolation{SrcLine: 39, Type: linter.CONDITION_EVALUATED_STATICALLY},
		actualViolation{SrcLine: 40, Type: linter.CONDITION_EVALUATED_STATICALLY},
	}

	if err := verifyViolations(expectedViolations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: NO_BUFFERED_FLUSHING
// Always flush the buffer when terminating, else the last output will not written.
func TestDetectionOfNotFlushingBufferedWriter(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_bufferwriting.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		actualViolation{SrcLine: 14, Type: linter.BUFFER_NOT_FLUSHED},
	}

	if err := verifyViolations(expectedViolations, actualViolations); err != nil {
		t.Fatal(err)
	}
}

*/

// Issue #4 list a scenario where the tool detects a false positive of rule ERROR_IGNORED.
// This test verifies correction of the behaviour.
func TestDetectionOfErrorIgnoredInReturn(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_errorinreturn.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{}

	if len(expectedViolations) > 0 {
		// Only verify violations of there is more than one file containing violations!
		if err := verifyViolations(expectedViolations[0].Violations, actualViolations); err != nil {
			t.Fatal(err)
		}
	}
}
