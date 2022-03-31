package property

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

type Analyzer struct {
	stCh chan *BiuStruct
}

func (a *Analyzer) Walk(f ast.Node) <-chan *BiuStruct {
	a.stCh = make(chan *BiuStruct, 1)
	go func() {
		defer close(a.stCh)
		ast.Walk(a, f)
	}()
	return a.stCh
}

func (a *Analyzer) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.TypeSpec:
		st := node.(*ast.TypeSpec)
		if structType, ok := st.Type.(*ast.StructType); ok {
			parent := &BiuField{}
			parent.Name = "i"
			parent.Type = st.Name.Name
			fields := make([]*BiuField, 0)
			for _, f := range structType.Fields.List {
				bField := &BiuField{}
				bField.Name = f.Names[0].Name
				if f.Comment != nil {
					cs := make([]string, 0)
					for _, c := range f.Comment.List {
						comment := strings.ReplaceAll(c.Text, "//", "")
						cs = append(cs, comment)
					}
					bField.Comments = cs
				}
				bField.Parent = parent
				switch f.Type.(type) {
				case *ast.Ident:
					bField.Type = f.Type.(*ast.Ident).Name
				case *ast.SelectorExpr:
					s := f.Type.(*ast.SelectorExpr)
					bField.Type = FlatSelectorExpr(s)
				}
				fields = append(fields, bField)
			}
			biuStruct := &BiuStruct{}
			biuStruct.Name = st.Name.Name
			biuStruct.Fields = fields
			a.stCh <- biuStruct
			// (&Generator{}).Invoke(biuStruct)
		}
	}
	return a
}

func FlatSelectorExpr(node *ast.SelectorExpr) string {
	var (
		prefix string
	)
	l := node.X
	r := node.Sel
	switch l.(type) {
	case *ast.Ident:
		prefix = l.(*ast.Ident).Name
	case *ast.SelectorExpr:
		prefix = FlatSelectorExpr(l.(*ast.SelectorExpr))
	}
	return fmt.Sprintf("%s.%s", prefix, r.Name)
}

func Input() error {
	var (
		packagePath string
	)
	packagePath = "D:\\WorkSpace\\gitee\\go_others\\property\\demo"

	fSet := token.NewFileSet()
	pkMap, err := parser.ParseDir(fSet, packagePath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		println(err)
	}
	for _, pk := range pkMap {
		p := &BiuPackage{}
		p.Name = pk.Name
		p.BiuFiles = make([]*BiuFile, 0)

		for k, f := range pk.Files {
			_, fileNameFmt := path.Split(filepath.ToSlash(k))
			if strings.Contains(fileNameFmt, "_prop.go") {
				continue
			}
			fileSuffix := filepath.Ext(fileNameFmt)
			fileName := strings.TrimSuffix(fileNameFmt, fileSuffix)

			biuFile := &BiuFile{}
			biuFile.Name = fmt.Sprintf("%s_prop.go", fileName)
			biuFile.Path = packagePath
			biuFile.Package = p

			v := &Analyzer{}
			for st := range v.Walk(f) {
				biuFile.BiuStructs = append(biuFile.BiuStructs, st)
			}
			// write 2 file
			(&Generator{}).Invoke(biuFile)
		}
	}
	return nil
}

func Visitor() {

}
