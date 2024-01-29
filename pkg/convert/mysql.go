package convert

import "fmt"

func ToMysqlType(gotype string, size int64) string {
	switch gotype {
	case "int8", "uint8":
		return "tinyint"
	case "int32", "uint32":
		return "int"
	case "int", "int64", "uint", "uint64":
		return "bigint"
	case "string":
		switch size {
		case -1:
			return "text"
		case -2:
			return "mediumtext"
		case -3:
			return "longtext"
		default:
			if size == 0 {
				size = 255
			}
			return fmt.Sprintf("varchar(%d)", size)
		}
	case "":
		return "embedded"
	default:
		return "embedded"
	}
}
