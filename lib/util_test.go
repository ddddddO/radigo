package lib

import (
	"testing"
)

func TestIsExist(t *testing.T) {
	if IsExist("XXXX.txt") {
		t.Error("failed")
	}

	if !IsExist("./util.go") {
		t.Error("failed")
	}
}