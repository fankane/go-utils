package file

import (
	"fmt"
	"testing"
)

func TestWriteCSV(t *testing.T) {
	lines := make([][]string, 0)
	lines = append(lines, []string{"1", "2,4", "3"})
	lines = append(lines, []string{"a", "ds"})
	if err := WriteCSV("./tt7.csv", lines); err != nil {
		fmt.Println("err:", err)
		return
	}

	lines2 := make([][]string, 0)
	lines2 = append(lines2, []string{"q1", "12,4", "13"})
	lines2 = append(lines2, []string{"1a", "1ds"})
	if err := WriteCSV("./tt7.csv", lines2); err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("success")
}
