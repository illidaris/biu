package convert

func ToMysqlType(gotype string) string {
	switch gotype {
	case "int8", "uint8":
		return "tinyint"
	case "int32", "uint32":
		return "int"
	case "int", "int64", "uint", "uint64":
		return "bigint"
	case "string":
		return "varchar"
	default:
		return ""
	}
}
