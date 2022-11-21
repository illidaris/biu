package analyzer

import "testing"

func TestRead(t *testing.T) {
	m, err := Visit("D:\\WorkSpace\\ztgame\\op\\norgannon\\adapter\\po")
	if err != nil {
		t.Error(err)
	}
	Read(m)
}
