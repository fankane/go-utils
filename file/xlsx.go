package file

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

/**
excelize 官方文档:https://xuri.me/excelize/zh-hans/base/installation.html#read
*/

const defaultSheetIdx = 1

func XlsxSheetMap(filePath string) (map[int]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.GetSheetMap(), nil
}

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
func XlsxAddColAtLast(filePath string, colVals []interface{}) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	sheetName, err := getSheetName(f, defaultSheetIdx)
	if err != nil {
		return err
	}
	if err = addColAtLast(f, sheetName, colVals); err != nil {
		return err
	}
	return nil
}

func XlsxSheetAddColAtLast(filePath, sheetName string, colVals []interface{}) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return addColAtLast(f, sheetName, colVals)
}

// XlsxAddRowsAtLast 在最后面插入一行
func XlsxAddRowsAtLast(filePath string, rowVals []interface{}) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	sheetName, err := getSheetName(f, defaultSheetIdx)
	if err != nil {
		return err
	}
	if err = addRowAtLast(f, sheetName, rowVals); err != nil {
		return err
	}
	return nil
}

// XlsxSheetAddRowsAtLast 在指定Sheet最后面插入一行
func XlsxSheetAddRowsAtLast(filePath, sheetName string, rowVals []interface{}) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = addRowAtLast(f, sheetName, rowVals); err != nil {
		return err
	}
	return nil
}

// DelRowsWithSheetName 删除指定行, start 从1开始
// end: -1 表示删除到最后一行
func DelRowsWithSheetName(filePath, sheetName string, start, end int) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return doDelRowsWithSheetName(f, sheetName, start, end)
}

// DelRowsWithSheetIdx 删除指定行, start 从1开始
// end: -1 表示删除到最后一行
func DelRowsWithSheetIdx(filePath string, sheetIdx, start, end int) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	sheetName, err := getSheetName(f, sheetIdx)
	if err != nil {
		return fmt.Errorf("get sheet name err:%s", err)
	}
	return doDelRowsWithSheetName(f, sheetName, start, end)
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

func addColAtLast(f *excelize.File, sheetName string, colVals []interface{}) error {
	if sheetName == "" {
		return fmt.Errorf("sheetName is empty")
	}
	colNum, err := getColNumOfRow(f, sheetName, 1)
	if err != nil {
		return err
	}
	/**
	insertCol 是在指定列的【前面】添加,所以需要+2
	比如：只有一列A，需要往后面添加B列，则需要在 C列前插入B， A-C的距离为2
	*/
	lastColName, err := excelize.ColumnNumberToName(colNum + 2)
	if err != nil {
		return err
	}
	if err = f.InsertCols(sheetName, lastColName, 1); err != nil {
		return err
	}
	if len(colVals) > 0 {
		newColName, err := excelize.ColumnNumberToName(colNum + 1)
		if err != nil {
			return err
		}
		if err = f.SetSheetCol(sheetName, fmt.Sprintf("%s1", newColName), &colVals); err != nil {
			return err
		}
	}

	if err := f.Save(); err != nil {
		return fmt.Errorf("file save err:%s", err)
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

// 获取行数
func getSheetRowLen(f *excelize.File, sheetName string) (int, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0, err
	}
	return len(rows), nil
}

func addRowAtLast(f *excelize.File, sheetName string, rowVals []interface{}) error {
	if sheetName == "" {
		return fmt.Errorf("sheetName is empty")
	}
	rowIdx, err := getSheetRowLen(f, sheetName)
	if err != nil {
		return err
	}
	/**
	insertRow 是在指定行的【前面】添加,所以需要+2
	比如：只有一行 1，需要往后面添加第2行，则需要在第3行前插入第2 行，
	*/
	if err = f.InsertRows(sheetName, rowIdx+2, 1); err != nil {
		return err
	}
	if len(rowVals) > 0 {
		if err = f.SetSheetRow(sheetName, fmt.Sprintf("A%d", rowIdx+1), &rowVals); err != nil {
			return err
		}
	}

	if err = f.Save(); err != nil {
		return fmt.Errorf("file save err:%s", err)
	}
	return nil
}

func doDelRowsWithSheetName(f *excelize.File, sheetName string, start, end int) error {
	if start < 1 {
		return fmt.Errorf("start must be greater than 0")
	}
	if end != -1 && start > end {
		return fmt.Errorf("start must be less than or equal end")
	}

	// 获取所有行数据
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("读取工作表失败:%s", err)
	}
	if len(rows) < start {
		return fmt.Errorf("行数不足, start:%d, rows:%d", start, len(rows))
	}
	if end == -1 {
		end = len(rows)
	}
	for i := end; i >= start; i-- {
		if err = f.RemoveRow(sheetName, i); err != nil {
			return fmt.Errorf("remove %d 行失败 err:%s", i, err)
		}
	}
	return f.Save()
}
