package file

import "testing"

func TestXlsxAddRowsAtLast(t *testing.T) {
	data := make([]interface{}, 0)
	data = append(data, "1")
	data = append(data, "2")
	data = append(data, "3")
	data = append(data, "4")
	data = append(data, "5")
	data = append(data, "6")
	data = append(data, "7")
	data = append(data, "8")
	data = append(data, "9")
	data = append(data, "10")
	data = append(data, "11")
	data = append(data, "12")
	XlsxAddRowsAtLast("./test.xlsx", data)
}
