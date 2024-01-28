package excel2model

import (
	"biu/internal/po"
	"biu/pkg/convert"
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/token"
	"log"
	"os"
	"path"
	"strings"

	"github.com/xuri/excelize/v2"
)

func Invoke(filepath string) {
	ch, _ := ReadFrmExcel(filepath)
	for v := range ch {
		f := v.ToAstFile()
		if f == nil {
			continue
		}
		Write2File(fmt.Sprintf("tmp/%s.go", v.GetAstName()), token.NewFileSet(), f)
	}
}

func ReadFrmExcel(filepath string) (<-chan *po.GoFile, error) {
	dir, filename := path.Split(filepath)
	println(dir)
	keys := strings.Split(filename, ".")
	if len(keys) == 0 {
		return nil, errors.New("filename is err")
	}
	ch := make(chan *po.GoFile)
	go func() {
		defer close(ch)
		f, err := excelize.OpenFile(filepath)
		if err != nil {
			println(err.Error())
		}
		defer func() {
			if err := f.Close(); err != nil {
				println(err.Error())
			}
		}()
		sheets := f.GetSheetList()
		for _, sheet := range sheets {
			rows, err := f.GetRows(sheet)
			if err != nil {
				println(err.Error())
				continue
			}
			sts := []*po.GoStruct{}
			st := &po.GoStruct{}
			st.Name = convert.Ucfirst(sheet) // 结构体名
			st.Fields = po.Rows2Struct[po.GoField](rows...)
			sts = append(sts, st)
			fl := &po.GoFile{}
			fl.Name = sheet          // 文件名
			fl.PackageName = keys[0] // 包名
			fl.Structs = sts
			ch <- fl
		}
	}()
	return ch, nil
}

func Write2File(dst string, set *token.FileSet, node interface{}) {
	var output []byte
	buffer := bytes.NewBuffer(output)
	err := format.Node(buffer, set, node)
	if err != nil {
		log.Fatal(err)
	}
	// 输出Go代码
	os.WriteFile(dst, buffer.Bytes(), os.ModePerm)
}

// wSet := token.NewFileSet()
// wf := &ast.File{}

// for _, st := range file.BiuStructs {
// 	for _, f := range st.Fields {
// 		wf.Decls = append(wf.Decls, i.SetterFunc(f))
// 		wf.Decls = append(wf.Decls, i.GetterFunc(f))
// 	}
// }

// wf.Name = &ast.Ident{
// 	Name: file.Package.Name,
// }

// Write2File(path.Join(file.Path, file.Name), wSet, wf)
