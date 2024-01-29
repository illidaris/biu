package convert

import (
	"reflect"
	"strconv"
	"strings"
)

func Struct2Row[T any](raw T) []string {
	rv := reflect.TypeOf(raw)
	row := Type2Row(rv)
	println(rv.NumField())
	return row
}

func Type2Row(t reflect.Type) []string {
	row := []string{}
	num := t.NumField()
	for i := 0; i < num; i++ {
		f := t.Field(i)
		switch f.Type.Kind() {
		case reflect.Struct:
			subRow := Type2Row(f.Type)
			row = append(row, subRow...)
			continue
		default:
			row = append(row, f.Name)
		}
	}
	return row
}

func Row2Struct[T any](rows ...[]string) []*T {
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
