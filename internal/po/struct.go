package po

import "go/ast"

type GoStruct struct {
	NameSection
	Fields []*GoField
}

func (s *GoStruct) ToAstStruct() *ast.TypeSpec {
	st := &ast.TypeSpec{}
	fields := []*ast.Field{}
	for _, f := range s.Fields {
		fields = append(fields, f.ToAstField())
	}
	st.Name = s.GetAstName()
	st.Type = &ast.StructType{
		Fields: &ast.FieldList{List: fields},
	}
	return st
}
