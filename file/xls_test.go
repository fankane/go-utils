package file

import (
	"fmt"
	"testing"
)

func TestReadXls(t *testing.T) {
	//m, err := XlsBaseInfo("./test3.xls")
	//if err != nil {
	//	fmt.Println("err:", err)
	//	return
	//}
	//fmt.Println(m.SheetNum, m.SheetMap[0])
	//return
	//fmt.Println("num:", num)
	//names, err := XlsSheetNames("./test.xls")
	//if err != nil {
	//	fmt.Println("err:", err)
	//	return
	//}
	//fmt.Println("names:", names)
	//res, err := ReadXlsWithSheetName("./test.xls", "Sheet1")
	res, err := XlsxSheet("E:\\app\\resource\\1688631230981752\\x-test.xlsx")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("行:", len(res))
	fmt.Println("列:", len(res[0]))
	for _, re := range res {
		for _, s := range re {
			fmt.Print(s, ",")
		}
		fmt.Println("")
	}
}
