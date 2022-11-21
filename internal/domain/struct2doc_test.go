package domain

import "testing"

// filePath:="D:\\WorkSpace\\ztgame\\op\\norgannon"
func TestInput(t *testing.T) {
	sd := &Struct2Doc{Source: "D:\\WorkSpace\\gitee\\go_others\\project\\po", Target: "./out"}
	err := sd.Invoke()
	if err != nil {
		t.Error(err)
	}
}
