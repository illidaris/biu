package po

import (
	"biu/pkg/convert"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

func Rows2Struct[T any](rows ...[]string) []*T {
	var data []*T
	// 默认第一行对应tag
	head := rows[0]
	for _, row := range rows[1:] {
		stu := new(T)
		rv := reflect.ValueOf(stu).Elem()
		for i := 0; i < len(row); i++ {
			colCell := row[i]
			// 通过 tag 取到结构体字段下标
			colCell = strings.Trim(colCell, " ")
			// 通过字段下标找到字段放射对象
			fName := head[i]
			v := rv.FieldByName(fName)
			// 根据字段的类型，选择适合的赋值方法
			switch v.Kind() {
			case reflect.String:
				value := colCell
				v.SetString(value)
			case reflect.Int64, reflect.Int32, reflect.Int8:
				value, _ := strconv.Atoi(colCell)
				// if err != nil {
				// 	panic(err)
				// }
				v.SetInt(int64(value))
			case reflect.Float64:
				value, _ := strconv.ParseFloat(colCell, 64)
				// if err != nil {
				// 	panic(err)
				// }
				v.SetFloat(value)
			}
		}

		data = append(data, stu)
	}
	return data
}

type GoField struct {
	NameSection
	TypeSection
	Size      int64  `json:"size"`
	Pk        int8   `json:"pk"`
	NotNull   int8   `json:"NotNull"`
	AutoIncr  int64  `json:"autoIncr"`
	OmitEmpty int8   `json:"omitempty"`
	Default   string `json:"default"`
	Unique    string `json:"unique"`
	Index     string `json:"index"`
	Enum      string `json:"enum"`
	Comment   string `json:"comment"`
	Extends   string `json:"extends"`
	Describe  string `json:"describe"`
}

func (i GoField) GetComment() *ast.CommentGroup {
	return &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Slash: token.Pos(token.DEFINE),
				Text:  fmt.Sprintf(" // %s %s", i.Comment, i.Describe),
			},
		},
	}
}

func (i GoField) GetJsonTag() string {
	if i.OmitEmpty == 0 {
		return fmt.Sprintf("json:\"%s\"",
			convert.Lcfirst(i.Name),
		)
	}
	return fmt.Sprintf("json:\"%s,omitempty\"",
		convert.Lcfirst(i.Name),
	)
}

func (i GoField) GetGormTag() string {
	keys := []string{
		"column:" + convert.Lcfirst(i.Name),
		"type:" + convert.ToMysqlType(i.Type),
	}
	if i.Size > 0 {
		keys = append(keys, fmt.Sprintf("size:%d", i.Size))
	}
	if i.Pk != 0 {
		keys = append(keys, "primaryKey")
	}
	if len(i.Comment) > 0 {
		keys = append(keys, "comment:"+i.Comment)
	}
	if len(i.Extends) > 0 {
		keys = append(keys, i.Extends)
	}
	return `gorm:"` + strings.Join(keys, ";") + `"`
}

func (i GoField) GetTag() *ast.BasicLit {
	return &ast.BasicLit{
		Kind: token.STRING,
		Value: fmt.Sprintf("`%s`", strings.Join([]string{
			i.GetJsonTag(),
			i.GetGormTag(),
		}, " ")),
	}
}

func (i GoField) ToAstField() *ast.Field {
	field := &ast.Field{}
	field.Names = i.GetAstNames()
	field.Type = i.GetAstType()
	field.Comment = i.GetComment()
	field.Tag = i.GetTag()
	return field
}
