package property

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Generator struct{}

func (i *Generator) SetterNameFunc() func(string) string {
	return func(s string) string {
		return fmt.Sprintf("Set%s", s)
	}
}

func (i *Generator) GetterNameFunc() func(string) string {
	return func(s string) string {
		return fmt.Sprintf("Get%s", s)
	}
}

func (i *Generator) Invoke(file *BiuFile) error {
	wSet := token.NewFileSet()
	wf := &ast.File{}

	for _, st := range file.BiuStructs {
		for _, f := range st.Fields {
			wf.Decls = append(wf.Decls, i.SetterFunc(f))
			wf.Decls = append(wf.Decls, i.GetterFunc(f))
		}
	}

	wf.Name = &ast.Ident{
		Name: file.Package.Name,
	}

	Write2File(path.Join(file.Path, file.Name), wSet, wf)
	return nil
}

func (i *Generator) SetterFunc(field *BiuField) ast.Decl {
	// params
	value := &BiuField{}
	value.Name = "value"
	value.Type = field.Type

	params := &ast.FieldList{
		List: []*ast.Field{
			value.ToAstField(false),
		},
	}
	// results
	errResult := &BiuField{}
	errResult.Type = "error"

	results := &ast.FieldList{
		List: []*ast.Field{
			errResult.ToAstField(false),
		},
	}
	// body
	body := []ast.Stmt{
		// before
		&ast.IfStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "err",
					},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "i",
							},
							Sel: &ast.Ident{
								Name: "Before",
							},
						},
					},
				},
			},
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "err",
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		// after
		&ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "i"},
					Sel: &ast.Ident{Name: "After"},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("\"%s\"", field.Name),
					},
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("\"%s\"", field.Nick),
					},
					&ast.Ident{
						Name: "value",
					},
				},
			},
		},

		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   field.Parent.GetIdentName(),
					Sel: field.GetIdentName(),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				value.GetIdentName(),
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.Ident{
					Name: "nil",
				},
			},
		},
	}

	setter := &ast.FuncDecl{
		Recv: field.GetAstReceiver(),
		Name: field.GetIdentName(i.SetterNameFunc()),
		Type: &ast.FuncType{
			Params:  params,
			Results: results,
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
	setter.Doc = &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: fmt.Sprintf("\n// %s %s,%s", setter.Name.Name, "setter func", field.OneComment()),
			},
		},
	}
	return setter
}

func (i *Generator) GetterFunc(field *BiuField) ast.Decl {
	// results
	errResult := &BiuField{}
	errResult.Type = field.Type

	results := &ast.FieldList{
		List: []*ast.Field{
			errResult.ToAstField(false),
		},
	}
	// body
	body := []ast.Stmt{
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.SelectorExpr{
					X:   field.Parent.GetIdentName(),
					Sel: field.GetIdentName(),
				},
			},
		},
	}

	getter := &ast.FuncDecl{
		Recv: field.GetAstReceiver(),
		Name: field.GetIdentName(i.GetterNameFunc()),
		Type: &ast.FuncType{
			Results: results,
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
	getter.Doc = &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: fmt.Sprintf("\n// %s %s,%s", getter.Name.Name, "getter func", field.OneComment()),
			},
		},
	}
	return getter
}

func Write2File(filename string, fSet *token.FileSet, node interface{}) {
	var output []byte
	buffer := bytes.NewBuffer(output)
	err := format.Node(buffer, fSet, node)
	if err != nil {
		log.Fatal(err)
	}
	// 输出Go代码
	ioutil.WriteFile(filename, buffer.Bytes(), os.ModePerm)
}
