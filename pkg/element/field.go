package element

import "go/ast"

type BiuField struct {
	Named
	Typed
	Commented
	Nick                  string
	WithBefore, WithAfter bool
	Parent                *BiuField
}

func (f *BiuField) ToAstField(isStar bool) *ast.Field {
	return &ast.Field{
		Names: f.GetIdentNames(),
		Type:  f.GetIdentType(isStar),
		Tag:   nil, // TODO: tag
	}
}

func (f *BiuField) GetAstReceiver() *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{
			f.Parent.ToAstField(true),
		},
	}
}
