package convert

import "fmt"

func ConvertNumberToLetter(num int) string {
	if num < 1 || num > 26 { // 只能处理A到Z之间的数字
		return ""
	}

	letter := 'A' + rune(num-1) // ASCII码表示大写字母从65开始，所以需要进行相应的计算
	return fmt.Sprintf("%c", letter)
}
