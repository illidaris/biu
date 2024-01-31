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
		astSt := st.ToPoStruct()
		wf.Decls = append(wf.Decls,
			// 结构体
			&ast.GenDecl{
				Tok:   token.TYPE,
				Specs: []ast.Spec{astSt},
			},
			TableNameFunc(st.Name),
		)
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

func (i *GoFile) ToMDFile() [][]string {
	rows := [][]string{}
	for _, st := range i.Structs {
		rows = append(rows,
			[]string{},
			[]string{
				st.Name,
				"类型",
				"描述",
				"备注",
			},
			[]string{
				"---",
				"---",
				"---",
				"---",
			})
		for _, field := range st.Fields {
			rows = append(rows, []string{
				field.Name,
				field.Type,
				field.Comment,
				field.Describe,
			})
		}
	}
	return rows
}
