package file

import (
	"fmt"
	"testing"
)

func TestWriteCSV(t *testing.T) {

	o := "Hello 中国"
	//chZ := strconv.QuoteToASCII(o)
	lines := make([][]string, 0)
	lines = append(lines, []string{"1", "2,4", "3"})
	lines = append(lines, []string{"a", o})
	if err := WriteCSV("./tt5.csv", lines); err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("success")
}
