package file

import (
	"fmt"
	"github.com/extrame/xls"
	"strings"
)

/**
xls 格式的文件比较古老，只能读取
需要写操作的，建议使用 xlsx, csv 等格式
*/
const xlsSuf = ".xls"

type XlsInfo struct {
	SheetNum int
	SheetMap map[int]SheetInfo //key:sheet下标，val:sheet详情
}
type SheetInfo struct {
	Name   string
	Idx    int //sheet下标
	Row    int
	MaxCol int //每行的列数不一样，返回最多的一个
}

func XlsBaseInfo(filePath string) (*XlsInfo, error) {
	wb, err := openXlsFile(filePath)
	if err != nil {
		return nil, err
	}
	sheetM := make(map[int]SheetInfo)
	for i := 0; i < wb.NumSheets(); i++ {
		info := SheetInfo{
			Name: wb.GetSheet(i).Name,
			Idx:  i,
		}
		if !isEmptySheet(wb.GetSheet(i)) {
			info.Row = int(wb.GetSheet(i).MaxRow) + 1 //maxRow是下标，从0开始的
			info.MaxCol = maxCol(wb.GetSheet(i))
		}
		sheetM[i] = info
	}
	return &XlsInfo{
		SheetNum: wb.NumSheets(),
		SheetMap: sheetM,
	}, nil
}

// ReadXls 读取 xls 第一个sheet的内容
func ReadXls(filePath string) ([][]string, error) {
	return ReadXlsWithSheetIdx(filePath, 0)
}

// ReadXlsWithSheetIdx 读取 xls 指定 sheet的内容
func ReadXlsWithSheetIdx(filePath string, sheetNum int) ([][]string, error) {
	wb, err := openXlsFile(filePath)
	if err != nil {
		return nil, err
	}
	return readXlsSheet(wb, sheetNum)
}

// ReadXlsWithSheetName 读取 xls 指定 sheet的内容
func ReadXlsWithSheetName(filePath, sheetName string) ([][]string, error) {
	wb, err := openXlsFile(filePath)
	if err != nil {
		return nil, err
	}
	for i := 0; i < wb.NumSheets(); i++ {
		if wb.GetSheet(i).Name == sheetName {
			return readXlsSheet(wb, i)
		}
	}
	return nil, fmt.Errorf("sheetName:%s not exist", sheetName)
}

// XlsSheetNum xls sheet数量
func XlsSheetNum(filePath string) (int, error) {
	wb, err := openXlsFile(filePath)
	if err != nil {
		return 0, err
	}
	return wb.NumSheets(), nil //NumSheets 是从第0行开始计算的，返回数量的时候，从1开始计算
}

// XlsSheetNames xls sheet名字列表
func XlsSheetNames(filePath string) ([]string, error) {
	wb, err := openXlsFile(filePath)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, wb.NumSheets())
	for i := 0; i < wb.NumSheets(); i++ {
		names = append(names, wb.GetSheet(i).Name)
	}
	return names, nil
}

func XlsSheetRows(filePath string, sheetIdx int) (int, error) {
	wb, err := openXlsFile(filePath)
	if err != nil {
		return 0, err
	}
	if sheetIdx >= wb.NumSheets() {
		return 0, fmt.Errorf("idx not exist")
	}
	if isEmptySheet(wb.GetSheet(sheetIdx)) {
		return 0, nil
	}
	return int(wb.GetSheet(sheetIdx).MaxRow) + 1, nil
}

func XlsSheetRowsWithName(filePath, sheetName string) (int, error) {
	idx, err := XlsSheetIdx(filePath, sheetName)
	if err != nil {
		return 0, err
	}
	return XlsSheetRows(filePath, idx)
}

func XlsSheetIdx(filePath string, sheetName string) (int, error) {
	wb, err := openXlsFile(filePath)
	if err != nil {
		return 0, err
	}
	for i := 0; i < wb.NumSheets(); i++ {
		if sheetName == wb.GetSheet(i).Name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("do not exist")
}

func openXlsFile(filePath string) (*xls.WorkBook, error) {
	if !strings.HasSuffix(filePath, xlsSuf) {
		return nil, fmt.Errorf("not *.xls file")
	}
	if !Exist(filePath) {
		return nil, fmt.Errorf("file not exist")
	}
	xlsFile, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, err
	}
	if xlsFile == nil {
		return nil, fmt.Errorf("xlsFile is nil")
	}
	return xlsFile, nil
}

func readXlsSheet(xlsWB *xls.WorkBook, sheetNum int) ([][]string, error) {
	sheet := xlsWB.GetSheet(sheetNum)
	if sheet == nil {
		return nil, fmt.Errorf("sheet [No.%d] not exist", sheetNum)
	}
	res := make([][]string, 0)
	for i := 0; i <= int(sheet.MaxRow); i++ {
		res = append(res, getRow(sheet.Row(i)))
	}
	return res, nil
}

func getRow(row *xls.Row) []string {
	if row == nil {
		return []string{}
	}
	data := make([]string, 0, row.LastCol())
	for i := 0; i < row.LastCol(); i++ {
		data = append(data, row.Col(i))
	}
	return data
}

func isEmptySheet(sheet *xls.WorkSheet) bool {
	return sheet.MaxRow == 0 && sheet.Row(0) == nil
}

func maxCol(sheet *xls.WorkSheet) int {
	max := 0
	for i := 0; i <= int(sheet.MaxRow); i++ {
		if sheet.Row(i) == nil {
			continue
		}
		if max < sheet.Row(i).LastCol() {
			max = sheet.Row(i).LastCol()
		}
	}
	return max
}
