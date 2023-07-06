package file

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

/**
excelize 官方文档:https://xuri.me/excelize/zh-hans/base/installation.html#read
*/

const defaultSheetIdx = 1

func XlsxSheet(filePath string) ([][]string, error) {
	return XlsxSheetOfIdx(filePath, defaultSheetIdx)
}

func XlsxSheetWithName(filePath, sheetName string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.GetRows(sheetName)
}

func XlsxSheetOfIdx(filePath string, idx int) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return sheetOfIdx(f, idx)
}

// XlsxAddColAtLast 在最后面插入一列，根据第一行的列数来确定最后一列
func XlsxAddColAtLast(filePath string, colVals interface{}) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	sheetName, err := getSheetName(f, defaultSheetIdx)
	if err != nil {
		return err
	}
	return addColAtLast(f, sheetName, colVals)
}

func XlsxSheetAddColAtLast(filePath, sheetName string, colVals interface{}) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return addColAtLast(f, sheetName, colVals)
}

func sheetOfIdx(f *excelize.File, idx int) ([][]string, error) {
	name, err := getSheetName(f, idx)
	if err != nil {
		return nil, err
	}
	content, err := f.GetRows(name)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func getSheetName(f *excelize.File, idx int) (string, error) {
	sm := f.GetSheetMap()
	if sm == nil || len(sm) == 0 {
		return "", fmt.Errorf("empty file")
	}
	name, ok := sm[idx]
	if !ok {
		return "", fmt.Errorf("no.[%d] not exist", idx)
	}
	return name, nil
}

func addColAtLast(f *excelize.File, sheetName string, colVals interface{}) error {
	if sheetName == "" {
		return fmt.Errorf("sheetName is empty")
	}
	colNum, err := getColNumOfRow(f, sheetName, 1)
	if err != nil {
		return err
	}
	fmt.Println("colNum:", colNum)
	lastColName, err := excelize.ColumnNumberToName(colNum)
	if err != nil {
		return err
	}
	fmt.Println("lastColName:", lastColName)
	if err = f.InsertCols(sheetName, lastColName, 1); err != nil {
		return err
	}

	newColName, err := excelize.ColumnNumberToName(colNum + 1)
	if err != nil {
		return err
	}
	startColName := fmt.Sprintf("%s1", newColName)
	fmt.Println("newColName:", startColName)
	//if err = f.SetSheetCol(sheetName, startColName, colVals); err != nil {
	//	return err
	//}
	if err = f.SetCellValue(sheetName, startColName, "你好吗"); err != nil {
		return err
	}

	return nil
}

// 获取 第n行的列数
func getColNumOfRow(f *excelize.File, sheetName string, n int) (int, error) {
	rows, err := f.Rows(sheetName)
	if err != nil {
		return 0, err
	}
	i := 0
	for rows.Next() {
		i++
		if i == n {
			cols, err := rows.Columns()
			if err != nil {
				return 0, err
			}
			return len(cols), nil
		}
	}
	return 0, fmt.Errorf("do not found no.[%d] row", n)
}
