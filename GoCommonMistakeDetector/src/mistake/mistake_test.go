package mistake

import (
	"testing"
	"io/ioutil"
)

func TestEquality(t *testing.T) {
	t.Error("Failed")
}

func TestUsageOfLoopVariableInGoroutine(t *testing.T) {
	srcFile, err := "/testSourceCode/loop_goroutine_mistake.go"
	if err != nil {
		t.Error(err)
	}

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Error(err)
	}

	_, err = FindCommonMistakes(sourceFile)
	if err != nil {
		t.Error("Failed to do something")
	}

}
