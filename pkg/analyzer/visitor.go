package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
)

func Visit(filePath string) (pkgs map[string]*ast.Package, first error) {
	fSet := token.NewFileSet()
	return parser.ParseDir(fSet, filePath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
}

func Read(pkgs map[string]*ast.Package) {
	for _, pkg := range pkgs {
		fmt.Printf("包名：%s", pkg.Name)
		fmt.Printf("包名：%s", pkg.Scope)
		fmt.Printf("包名：%v", pkg.Imports)
		fmt.Printf("包名：%v", pkg.Files)
	}
	fmt.Println("end")
}
