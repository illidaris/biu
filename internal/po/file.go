package po

import (
	"go/ast"
	"go/token"
)

type GoFile struct {
	NameSection
	Structs     []*GoStruct
	PackageName string
}

// GetPackageName 获取字段名声明
func (i *GoFile) GetPackageName(formats ...func(string) string) *ast.Ident {
	if len(i.PackageName) == 0 {
		return nil
	}
	name := i.PackageName
	for _, f := range formats {
		name = f(name)
	}
	return &ast.Ident{
		Name: name,
	}
}

func (i *GoFile) ToPoFile() *ast.File {
	wf := &ast.File{
		Decls: []ast.Decl{},
	}
	wf.Name = i.GetPackageName()
	for _, st := range i.Structs {
		wf.Decls = append(wf.Decls, &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{st.ToPoStruct()},
		})
	}
	return wf
}

func (i *GoFile) ToDtoFile() *ast.File {
	wf := &ast.File{
		Decls: []ast.Decl{},
	}
	wf.Name = i.GetPackageName()
	for _, st := range i.Structs {
		wf.Decls = append(wf.Decls, &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{st.ToDtoStruct()},
		})
	}
	return wf
}
