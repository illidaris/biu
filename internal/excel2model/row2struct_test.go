package excel2model

import (
	"biu/internal/po"
	"biu/pkg/convert"
	"testing"
)

func TestRow2Struct(t *testing.T) {
	convert.Struct2Row(po.GoField{})
}
