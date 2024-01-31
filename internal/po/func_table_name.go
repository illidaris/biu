package po

import (
	"biu/pkg/convert"
	"fmt"
	"go/ast"
	"go/token"
)

// TableNameFunc 获取TableName函数
func TableNameFunc(stName string) *ast.FuncDecl {
	// 函数
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						&ast.Ident{
							Name: "i",
						},
					},
					Type: &ast.Ident{
						Name: stName,
					},
				},
			}, // 接收
		},
		Name: &ast.Ident{
			Name: "TableName",
		},
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "string",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf("\"%s\"", convert.CasedName(stName)),
						},
					},
				},
			},
		},
	}
}
