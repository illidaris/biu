package excel2model

import (
	"biu/internal/po"
	"biu/pkg/convert"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"os"
	"path"
	"strings"

	fileex "github.com/illidaris/file/path"
	"github.com/xuri/excelize/v2"
)

func GenTemplate() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Set value of a cell.
	for index, val := range convert.Struct2Row(po.GoField{}) {
		f.SetCellValue("Sheet1", fmt.Sprintf("%s1", convert.ConvertNumberToLetter(index+1)), val)
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs("rename_package.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func Invoke(filepath string) {
	ch, _ := ReadFrmExcel(filepath)
	for v := range ch {
		po := v.ToPoFile()
		Write2File("tmp", "po", v.PackageName, v.Name, po)
		dto := v.ToDtoFile()
		Write2File("tmp", "dto", v.PackageName, v.Name, dto)
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
			st.Fields = convert.Row2Struct[po.GoField](rows...)
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

func Write2File(out string, module, pkg, name string, gofile *ast.File) {
	if gofile == nil {
		return
	}
	var set = &token.FileSet{}
	var output []byte
	var keys = []string{
		out,
		module,
		pkg,
	}
	_ = fileex.MkdirIfNotExist(path.Join(keys...))
	buffer := bytes.NewBuffer(output)
	err := format.Node(buffer, set, gofile)
	if err != nil {
		log.Fatal(err)
	}
	keys = append(keys, fmt.Sprintf("%s.go", name))
	// 输出Go代码
	os.WriteFile(path.Join(keys...), buffer.Bytes(), os.ModePerm)
}
