// struct => doc
package domain

import (
	"biu/pkg/analyzer"
	"bytes"
	"fmt"
	"go/ast"
	"io/fs"
	"io/ioutil"
	"path"
	"strings"

	pathEx "github.com/illidaris/file/path"
)

type Struct2Doc struct {
	Source, Target string
}

func (s *Struct2Doc) Invoke() error {
	if err := pathEx.MkdirIfNotExist(s.Target); err != nil {
		return err
	}
	pkgs, err := analyzer.Visit(s.Source)
	if err != nil {
		return err
	}

	// 输出根路径
	//outputDir := DirJion("E:", "output")
	for _, pkg := range pkgs {
		if pkg == nil {
			continue
		}
		a := &analyzer.Analyzer{
			F: s.Visit,
		}
		ast.Walk(a, pkg)
		// 解析包
		// pkgDir := DirJion(outputDir, pkg.Name)
		// fs := pkg.Files
	}

	return nil
}

func (s *Struct2Doc) Visit(node ast.Node) {
	switch node.(type) {
	case *ast.TypeSpec:
		ts := node.(*ast.TypeSpec)
		s.TypeSpec2Doc(ts)
	default:
		//println(node)
	}

}

func (s *Struct2Doc) TypeSpec2Fmt(ts *ast.TypeSpec) {
	structName := ts.Name.String()
	fmt.Printf("#### %s \n", structName)
	st := ts.Type
	switch st := st.(type) {
	case *ast.StructType:
		fmt.Println("|字段|字段类型|描述|标签|")
		fmt.Println("|---|---|---|---|")
		// 字段列表
		for _, v := range st.Fields.List {
			var (
				fieldName    = "-"
				fieldType    = "-"
				fieldComment = "-"
				fieldTag     = "-"
			)
			if len(v.Names) > 0 {
				fieldName = v.Names[0].String()
			}
			if v.Type != nil {
				fieldType = ParseExpr(v.Type)
			}
			if v.Comment != nil {
				fieldComment = v.Comment.Text()
			}
			if v.Tag != nil {
				fieldTag = v.Tag.Value
			}
			fmt.Printf("|%s|%s|%s|%s| \n", fieldName, fieldType, fieldComment, fieldTag)
		}
	case *ast.FuncType:
		return
	case *ast.InterfaceType:
		return
	default:
		return
	}
}

func (s *Struct2Doc) TypeSpec2Doc(ts *ast.TypeSpec) {
	doc := &bytes.Buffer{}
	structName := ts.Name.String()
	doc.WriteString(fmt.Sprintf("#### %s\n", structName))
	st := ts.Type
	switch st := st.(type) {
	case *ast.StructType:
		doc.WriteString("|字段|字段类型|描述|标签|\n")
		doc.WriteString("|---|---|---|---|\n")
		// 字段列表
		for _, v := range st.Fields.List {
			var (
				fieldName    = "-"
				fieldType    = "-"
				fieldComment = "-"
				fieldTag     = "-"
			)
			if len(v.Names) > 0 {
				fieldName = v.Names[0].String()
			}
			if v.Type != nil {
				fieldType = ParseExpr(v.Type)
			}
			if v.Comment != nil {
				fieldComment = strings.ReplaceAll(v.Comment.Text(), "\n", "。")
			}
			if v.Tag != nil {
				fieldTag = v.Tag.Value
			}
			doc.WriteString(fmt.Sprintf("|%s|%s|%s|%s|\n", fieldName, fieldType, fieldComment, fieldTag))
		}
	case *ast.FuncType:
		return
	case *ast.InterfaceType:
		return
	default:
		return
	}
	ioutil.WriteFile(fmt.Sprintf("%s/%s.md", s.Target, structName), doc.Bytes(), fs.ModePerm)
}

func ParseExpr(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.String()
	case *ast.SelectorExpr:
		l := ParseExpr(expr.X)
		r := expr.Sel.String()
		return fmt.Sprintf("%s.%s", l, r)
	default:
		return "-"
	}
}

// DirJion 拼接并且生成目录
func DirJion(dir, sub string) string {
	res := path.Join(dir, sub)
	err := pathEx.MkdirIfNotExist(res)
	if err != nil {
		return ""
	}
	return res
}
