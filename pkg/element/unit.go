package element

import (
	"go/ast"
	"strings"
)

type Named struct {
	Name string
}

// GetIdentName 获取字段名声明
func (i *Named) GetIdentName(formats ...func(string) string) *ast.Ident {
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

func (i *Named) GetIdentNames() []*ast.Ident {
	if name := i.GetIdentName(); name != nil {
		return []*ast.Ident{name}
	}
	return nil
}

type Typed struct {
	Type string
}

func (i *Typed) GetIdentType(isStar bool) ast.Expr {
	if len(i.Type) == 0 {
		return nil
	}
	ident := &ast.Ident{
		Name: i.Type,
	}
	if !isStar {
		return ident
	}
	return &ast.StarExpr{
		X: ident,
	}
}

type Commented struct {
	Comments []string
}

func (i *Commented) OneComment() string {
	return strings.Join(i.Comments, ",")
}
