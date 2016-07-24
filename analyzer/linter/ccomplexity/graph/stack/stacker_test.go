package stack_test

import (
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity/graph/stack"
	"testing"
)

func TestStackSize(t *testing.T) {
	stacker := stack.Stack{}

	stacker.Push("First")
	stacker.Push("Seconds")

	if stacker.Len() != 2 {
		t.Errorf("Size of stack should be 2, not %d\n", stacker.Len())
	}

	stacker.Pop()

	if stacker.Len() != 1 {
		t.Errorf("Size of stack should be 1, not %d\n", stacker.Len())
	}
}
