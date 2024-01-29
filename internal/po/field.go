package po

import (
	"biu/pkg/convert"
	"fmt"
	"go/ast"
	"go/token"

	"strings"
)

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
		"type:" + convert.ToMysqlType(i.Type, i.Size),
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

func (i GoField) ToPoField() *ast.Field {
	field := &ast.Field{}
	field.Names = i.GetAstNames()
	field.Type = i.GetAstType()
	field.Comment = i.GetComment()
	field.Tag = i.GetTag()
	return field
}

func (i GoField) ToDtoField() *ast.Field {
	field := &ast.Field{}
	field.Names = i.GetAstNames()
	field.Type = i.GetAstType()
	field.Comment = i.GetComment()
	field.Tag = &ast.BasicLit{
		Kind: token.STRING,
		Value: fmt.Sprintf("`%s`", strings.Join([]string{
			i.GetJsonTag(),
		}, " ")),
	}
	return field
}
