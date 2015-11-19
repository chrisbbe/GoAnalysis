package mistake

import "testing"

func TestUsageOfLoopVariableInGoroutine(t *testing.T) {
	a := 1
	b := 1
	if a != b {
		t.Errorf("%d != %d", a, b)
	}
}
