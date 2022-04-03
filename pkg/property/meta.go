package property

import (
	"fmt"
	"go/ast"
	"go/token"
)

func BuildBeforeFunc() ast.Stmt {
	return &ast.IfStmt{
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
	}
}

// after
func BuildAfterFunc(name, nick string) ast.Stmt {
	return &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "i"},
				Sel: &ast.Ident{Name: "After"},
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", name),
				},
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", nick),
				},
				&ast.Ident{
					Name: "value",
				},
			},
		},
	}
}
