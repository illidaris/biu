package po

import "go/ast"

type GoStruct struct {
	NameSection
	Fields []*GoField
}

func (s *GoStruct) ToPoStruct() *ast.TypeSpec {
	st := &ast.TypeSpec{}
	fields := []*ast.Field{}
	for _, f := range s.Fields {
		fields = append(fields, f.ToPoField())
	}
	st.Name = s.GetAstName()
	st.Type = &ast.StructType{
		Fields: &ast.FieldList{List: fields},
	}
	return st
}
func (s *GoStruct) ToDtoStruct() *ast.TypeSpec {
	st := &ast.TypeSpec{}
	fields := []*ast.Field{}
	for _, f := range s.Fields {
		fields = append(fields, f.ToDtoField())
	}
	st.Name = s.GetAstName()
	st.Type = &ast.StructType{
		Fields: &ast.FieldList{List: fields},
	}
	return st
}
