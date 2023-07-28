package file

import (
	"encoding/csv"
	"fmt"
	"os"
)

func WriteCSV(filePath string, lines [][]string) error {
	var (
		csvFile *os.File
		err     error
	)
	defer csvFile.Close()
	if !Exist(filePath) {
		csvFile, err = os.Create(filePath)
	} else {
		csvFile, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	}
	if err != nil {
		return fmt.Errorf("open file [%s] err:%s", filePath, err)
	}

	w := csv.NewWriter(csvFile)
	if err = w.WriteAll(lines); err != nil {
		return fmt.Errorf("write err:%s", err)
	}
	w.Flush()
	return nil
}

func ReadCSV(filePath string) ([][]string, error) {
	csvFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("open file [%s] err:%s", filePath, err)
	}
	return csv.NewReader(csvFile).ReadAll()
}
