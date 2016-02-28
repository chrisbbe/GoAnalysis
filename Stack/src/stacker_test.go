package stack

import "testing"

func TestLenOfStack(t *testing.T) {

	s := Stack{}

	s.Push("First")
	s.Push("Seconds")

	if s.Len() != 2 {
		t.Errorf("Size of stack should be 2, not %d\n", s.Len())
	}

	s.Pop()

	if s.Len() != 1 {
		t.Errorf("Size of stack should be 1, not %d\n", s.Len())
	}

}
