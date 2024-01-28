package po

import (
	"go/ast"
	"strings"
)

type NameSection struct {
	Name string `json:"name"`
}

// GetIdentName 获取字段名声明
func (i *NameSection) GetAstName(formats ...func(string) string) *ast.Ident {
	if len(i.Name) == 0 {
		return nil
	}
	name := i.Name
	for _, f := range formats {
		name = f(name)
	}
	return &ast.Ident{
		Name: name,
	}
}

func (i *NameSection) GetAstNames() []*ast.Ident {
	if name := i.GetAstName(); name != nil {
		return []*ast.Ident{name}
	}
	return nil
}

type TypeSection struct {
	Ptr  bool   `json:"ptr"`
	Type string `json:"type"`
}

func (i *TypeSection) GetAstType() ast.Expr {
	if len(i.Type) == 0 {
		return nil
	}
	ident := &ast.Ident{
		Name: i.Type,
	}
	if !i.Ptr {
		return ident
	}
	return &ast.StarExpr{
		X: ident,
	}
}

type CommentSection struct {
	Comments []string `json:"comments"`
}

func (i *CommentSection) OneComment() string {
	return strings.Join(i.Comments, ",")
}
