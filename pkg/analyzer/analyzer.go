package analyzer

import (
	"biu/pkg/element"
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
	stCh chan *element.BiuStruct
	// typeSpecCh chan *ast.TypeSpec
	F func(node ast.Node)
}

func (a *Analyzer) Walk(f ast.Node) <-chan *element.BiuStruct {
	a.stCh = make(chan *element.BiuStruct, 1)
	go func() {
		defer close(a.stCh)
		ast.Walk(a, f)
	}()
	return a.stCh
}

func (a *Analyzer) Visit(node ast.Node) ast.Visitor {
	a.F(node)
	return a
	// switch node.(type) {
	// case *ast.TypeSpec:
	// 	st := node.(*ast.TypeSpec)
	// 	if structType, ok := st.Type.(*ast.StructType); ok {
	// 		parent := &element.BiuField{}
	// 		parent.Name = "i"
	// 		parent.Type = st.Name.Name
	// 		fields := make([]*element.BiuField, 0)
	// 		for _, f := range structType.Fields.List {
	// 			bField := &element.BiuField{}
	// 			if f.Names == nil {
	// 				continue
	// 			}
	// 			bField.Name = f.Names[0].Name
	// 			bField.Nick = format.CasedName(bField.Name)
	// 			if f.Comment != nil {
	// 				cs := make([]string, 0)
	// 				for _, c := range f.Comment.List {
	// 					comment := strings.ReplaceAll(c.Text, "//", "")
	// 					cs = append(cs, comment)
	// 				}
	// 				bField.Comments = cs
	// 			}
	// 			bField.Parent = parent
	// 			switch f.Type.(type) {
	// 			case *ast.Ident:
	// 				bField.Type = f.Type.(*ast.Ident).Name
	// 			case *ast.SelectorExpr:
	// 				s := f.Type.(*ast.SelectorExpr)
	// 				bField.Type = FlatSelectorExpr(s)
	// 			}
	// 			fields = append(fields, bField)
	// 		}
	// 		biuStruct := &element.BiuStruct{}
	// 		biuStruct.Name = st.Name.Name
	// 		biuStruct.Fields = fields
	// 		a.stCh <- biuStruct
	// 	}
	// case *ast.FuncDecl:
	// 	funcDecl := node.(*ast.FuncDecl)
	// 	if funcDecl.Recv != nil && funcDecl.Recv.List != nil {
	// 		for _, v := range funcDecl.Recv.List {
	// 			if v.Type != nil {
	// 				if sType, ok := (v.Type).(*ast.StarExpr); ok {
	// 					if idname, ok := sType.X.(*ast.Ident); ok {
	// 						bStruct := &element.BiuStruct{}
	// 						bStruct.Name = idname.Name
	// 						if funcDecl.Name.Name == "Before" {
	// 							bStruct.WithBefore = true
	// 							a.stCh <- bStruct
	// 						}
	// 						if funcDecl.Name.Name == "After" {
	// 							bStruct.WithAfter = true
	// 							a.stCh <- bStruct
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// return a
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

func Input(filePath string) error {
	fSet := token.NewFileSet()
	pkMap, err := parser.ParseDir(fSet, filePath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)

	if err != nil {
		println(err)
	}
	for _, pk := range pkMap {
		p := &element.BiuPackage{}
		p.Name = pk.Name
		p.BiuFiles = make([]*element.BiuFile, 0)

		for k, f := range pk.Files {
			_, fileNameFmt := path.Split(filepath.ToSlash(k))
			if strings.Contains(fileNameFmt, "_prop.go") {
				continue
			}
			fileSuffix := filepath.Ext(fileNameFmt)
			fileName := strings.TrimSuffix(fileNameFmt, fileSuffix)

			biuFile := &element.BiuFile{}
			biuFile.Name = fmt.Sprintf("%s_prop.go", fileName)
			biuFile.Path = filePath
			biuFile.Package = p

			beforeMap := make(map[string]bool)
			afterMap := make(map[string]bool)

			v := &Analyzer{}
			for st := range v.Walk(f) {
				if len(st.Fields) > 0 {
					biuFile.BiuStructs = append(biuFile.BiuStructs, st)
				}
				if st.WithBefore {
					beforeMap[st.Name] = st.WithBefore
				}
				if st.WithAfter {
					afterMap[st.Name] = st.WithAfter
				}
			}
			for _, st := range biuFile.BiuStructs {
				if v, ok := beforeMap[st.Name]; ok {
					st.WithBefore = v
				}
				if v, ok := afterMap[st.Name]; ok {
					st.WithAfter = v
				}
				for _, f := range st.Fields {
					f.WithBefore = st.WithBefore
					f.WithAfter = st.WithAfter
				}
			}
		}
	}
	return nil
}
