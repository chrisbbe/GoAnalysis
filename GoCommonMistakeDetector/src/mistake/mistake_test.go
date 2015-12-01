package mistake

import (
	"testing"
)

func TestEquality(t *testing.T) {

	if 1 == 2 {
		t.Error("Failed")
	}
}
