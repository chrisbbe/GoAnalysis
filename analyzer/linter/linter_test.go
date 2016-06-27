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
package linter_test

import (
	"fmt"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter"
	"testing"
)

type actualViolation struct {
	SrcLine int
	Type    linter.Rule
}

// verifyViolations checks the list of expected mistakes with the list of actual mistakes.
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
	expectedMistakes, err := linter.DetectViolations("./testcode/_fmtprinting.go")
	if err != nil {
		t.Fatal(err)
	}

	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 27, Type: linter.FMT_PRINTING},
		actualViolation{SrcLine: 28, Type: linter.FMT_PRINTING},
		actualViolation{SrcLine: 29, Type: linter.FMT_PRINTING},
	}

	if err := verifyViolations(expectedMistakes, actualMistakes); err != nil {
		t.Fatal(err)
	}
}

// TODO
// Testing rule: STRING_DEFINES_ITSELF
func TestDetectionOfFmtPrintingInStringMethod(t *testing.T) {
	expectedMistakes, err := linter.DetectViolations("./testcode/_stringMethod.go")
	if err != nil {
		t.Fatal(err)
	}

	correctMistakes := []actualViolation{
		actualViolation{SrcLine: 40, Type: linter.STRING_DEFINES_ITSELF},
	}

	if err := verifyViolations(expectedMistakes, correctMistakes); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: NEVER_ALLOCATED_MAP_WITH_NEW
// Allocating maps with new returns a nil pointer, therefor one should use make.
func TestDetectionOfAllocatingMapWithNew(t *testing.T) {
	expectedMistakes, _ := linter.DetectViolations("./testcode/_newmap.go")
	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 25, Type: linter.NEVER_ALLOCATED_MAP_WITH_NEW},
	}

	if len(expectedMistakes) != len(actualMistakes) {
		t.Fatalf("Number of mistakes should be %d, not %d!\n", len(actualMistakes), len(expectedMistakes))
	}

	for index, expectedMistake := range expectedMistakes {
		if actualMistakes[index].Type != expectedMistake.Type {
			t.Errorf("Error should be of type %s, and not type %s!\n", actualMistakes[index].Type, expectedMistake.Type)
		}
		if actualMistakes[index].SrcLine != expectedMistake.SrcLine {
			t.Errorf("Error should be found on line %d, and not on line %d!\n", actualMistakes[index].SrcLine, expectedMistake.SrcLine)
		}
	}

}

/*
// TODO: Fix this check!

// Testing rule: ERROR_IGNORED
// Errors should never be ignored.
func TestDetectionOfIgnoredErrors(t *testing.T) {
	expectedViolations, err := linter.DetectViolations("./testcode/_ignoreerror.go")
	if err != nil {
		t.Fatal(err)
	}

	actualViolations := []actualViolation{
		actualViolation{SrcLine: 33, Type: linter.ERROR_IGNORED},
		actualViolation{SrcLine: 38, Type: linter.ERROR_IGNORED},
		actualViolation{SrcLine: 43, Type: linter.ERROR_IGNORED},
	}

	if err := verifyViolations(expectedViolations, actualViolations); err != nil {
		t.Fatal(err)
	}
}
*/

// Testing rule: RACE_CONDITION
// Races will occur, since multiple Go-routines will share the same counter variable.
func TestDetectionOfRacesInLoopClosures(t *testing.T) {
	expectedMistakes, err := linter.DetectViolations("./testcode/_threadLooping.go")
	if err != nil {
		t.Fatal(err)
	}

	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 21, Type: linter.RACE_CONDITION},
	}

	if err := verifyViolations(expectedMistakes, actualMistakes); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: EMPTY_IF_BODY
// Empty if-bodies are unnecessary and ineffective.
func TestDetectionOfEmptyIfBody(t *testing.T) {
	expectedMistakes, err := linter.DetectViolations("./testcode/_emptyifbody.go")
	if err != nil {
		t.Fatal(err)
	}

	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 27, Type: linter.EMPTY_IF_BODY},
	}

	if err := verifyViolations(expectedMistakes, actualMistakes); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: EMPTY_ELSE_BODY
// Empty else-bodies are unnecessary and ineffective.
func TestDetectionOfEmptyElseBody(t *testing.T) {
	expectedMistakes, err := linter.DetectViolations("./testcode/_emptyelsebody.go")
	if err != nil {
		t.Fatal(err)
	}

	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 31, Type: linter.EMPTY_ELSE_BODY},
	}

	if err := verifyViolations(expectedMistakes, actualMistakes); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: EMPTY_FOR_BODY
// Empty for-bodies are unnecessary and ineffective.
func TestDetectionOfEmptyForBody(t *testing.T) {
	expectedMistakes, err := linter.DetectViolations("./testcode/_emptyforbody.go")
	if err != nil {
		t.Fatal(err)
	}

	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 26, Type: linter.EMPTY_FOR_BODY},
	}

	if err := verifyViolations(expectedMistakes, actualMistakes); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: GOTO_USED
// Jumping around in the code using GOTO (including BREAK, CONTINUE, GOTO, FALLTHROUGH)
// is considered confusing and harmful.
func TestDetectionOfGoTo(t *testing.T) {
	expectedMistakes, err := linter.DetectViolations("./testcode/_goto.go")
	if err != nil {
		t.Fatal(err)
	}

	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 30, Type: linter.GOTO_USED},
	}

	if err := verifyViolations(expectedMistakes, actualMistakes); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: GOTO_USED
// Jumping around in the code using GOTO (including BREAK, CONTINUE, GOTO, FALLTHROUGH)
// is considered confusing and harmful.
func TestDetectionOfLabeledBranching(t *testing.T) {
	expectedMistakes, err := linter.DetectViolations("./testcode/_labeledbranch.go")
	if err != nil {
		t.Fatal(err)
	}

	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 33, Type: linter.GOTO_USED},
	}

	if err := verifyViolations(expectedMistakes, actualMistakes); err != nil {
		t.Fatal(err)
	}
}

/*

// TODO: Implement checks for this two rules!
// Testing rule: CONDITION_EVALUATED_STATICALLY
// Condition that can be evaluated statically are wasted and performance-reducing.
func TestDetectionOfConditionEvaluatedStatically(t *testing.T) {
	expectedMistakes, err := linter.DetectViolations("./testcode/_staticconditions.go")
	if err != nil {
		t.Fatal(err)
	}

	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 27, Type: linter.CONDITION_EVALUATED_STATICALLY},
		actualViolation{SrcLine: 31, Type: linter.CONDITION_EVALUATED_STATICALLY},
		actualViolation{SrcLine: 35, Type: linter.CONDITION_EVALUATED_STATICALLY},
	}

	if err := verifyViolations(expectedMistakes, actualMistakes); err != nil {
		t.Fatal(err)
	}
}

// Testing rule: NO_BUFFERED_FLUSHING
// Always flush the buffer when terminating, else the last output will not written.
func TestDetectionOfNotFlushingBufferedWriter(t *testing.T) {
	expectedMistakes, err := linter.DetectViolations("./testcode/_bufwriting.go")
	if err != nil {
		t.Fatal(err)
	}

	actualMistakes := []actualViolation{
		actualViolation{SrcLine: 31, Type: linter.CONDITION_EVALUATED_STATICALLY},
	}

	if err := verifyViolations(expectedMistakes, actualMistakes); err != nil {
		t.Fatal(err)
	}
}
*/
