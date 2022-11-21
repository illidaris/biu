package format

import (
	"testing"
)

func TestCasedName(t *testing.T) {
	result := CasedName("AbcAdb")
	if result != "abc_adb" {
		t.Error("error")
	}
}
