package file

import (
	"bufio"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
)

// DirExist 判定目录是否存在
func DirExist(dir string) bool {
	fInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return fInfo.IsDir()
}

func Exist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else if err == nil {
		return true
	}
	return false
}

func DirFiles(dir string) ([]string, error) {
	dirEntryArr, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	fileNames := make([]string, 0, len(dirEntryArr))
	for _, temp := range dirEntryArr {
		if !temp.IsDir() {
			fileNames = append(fileNames, temp.Name())
		}
	}
	return fileNames, nil
}

func Content(fileHeader *multipart.FileHeader) ([]byte, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func ReadLine(filePath string) ([]string, error) {
	if !Exist(filePath) {
		return nil, os.ErrNotExist
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	r := bufio.NewReader(f)
	for {
		bytes, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return lines, err
		}
		lines = append(lines, string(bytes))
	}
	return lines, nil
}
