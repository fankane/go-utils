package file

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func DeleteFiles(filePath ...string) error {
	delFailed := make([]string, 0)
	for _, s := range filePath {
		if !FileExist(s) {
			continue //文件不存在，直接返回
		}
		if err := os.Remove(s); err != nil {
			delFailed = append(delFailed, s)
		}
	}
	if len(delFailed) == 0 {
		return nil
	}
	return fmt.Errorf(fmt.Sprintf("delete failed:[%s]", strings.Join(delFailed, ",")))
}

func DeleteDir(dir string) error {
	if dir == "" {
		return nil
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}

// DeleteDirFilesWithPref 删除 dir 目录下所有以 pref 开头的文件
func DeleteDirFilesWithPref(dir, pref string) error {
	prefFiles := make([]string, 0)
	fs, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if !f.IsDir() && strings.HasPrefix(f.Name(), pref) {
			prefFiles = append(prefFiles, f.Name())
		}
	}
	delFailed := make([]string, 0)
	for _, s := range prefFiles {
		f := filepath.Join(dir, s)
		if !FileExist(f) {
			continue //文件不存在，直接返回
		}
		if err = os.Remove(f); err != nil {
			delFailed = append(delFailed, f)
		}
	}
	if len(delFailed) == 0 {
		return nil
	}
	return fmt.Errorf(fmt.Sprintf("delete failed:[%s]", strings.Join(delFailed, ",")))
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

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else if err == nil {
		return true
	}
	return false
}

func DirExist(dir string) bool {
	fInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return fInfo.IsDir()
}

func Mkdir(dir string) error {
	if DirExist(dir) {
		return nil
	}
	err := os.MkdirAll(dir, 0755)       // 创建目录，权限为 0755
	if err != nil && !os.IsExist(err) { // 如果出现错误且不是目录已存在的情况，则返回错误信息
		return err
	}
	return nil
}
