package file

import (
	"fmt"
	"testing"
)

func TestXlsxAddColAtLast(t *testing.T) {
	col := make([]interface{}, 0)
	col = append(col, "结果")
	col = append(col, "world")
	if err := XlsxAddColAtLast("./x-test.xlsx", &col); err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("success")
}
