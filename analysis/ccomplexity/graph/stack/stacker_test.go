package stack_test

import (
	"github.com/chrisbbe/GoAnalysis/analysis/ccomplexity/graph/stack"
	"testing"
)

func TestStackSize(t *testing.T) {
	stacker := stack.Stack{}

	stacker.Push("First")
	stacker.Push("Seconds")

	if stacker.Len() != 2 {
		t.Errorf("Size of ccomplexity.graph.stack should be 2, not %d\n", stacker.Len())
	}

	stacker.Pop()

	if stacker.Len() != 1 {
		t.Errorf("Size of ccomplexity.graph.stack should be 1, not %d\n", stacker.Len())
	}
}
