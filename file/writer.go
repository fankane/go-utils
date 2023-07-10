package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// MkdirWhenNo 目录不存在则创建，存在则忽略
func MkdirWhenNo(dir string, perm os.FileMode) error {
	if DirExist(dir) {
		return nil
	}
	err := os.MkdirAll(dir, perm)       // 创建目录，权限为 0755
	if err != nil && !os.IsExist(err) { // 如果出现错误且不是目录已存在的情况，则返回错误信息
		return err
	}
	return nil
}

// MkdirY 在dir下面，按年创建文件夹
func MkdirY(dir string, perm os.FileMode) error {
	now := time.Now()
	newDir := filepath.Join(dir, strconv.Itoa(now.Year()))
	return MkdirWhenNo(newDir, perm)
}

// MkdirYM 在dir下面，按年月创建文件夹
func MkdirYM(dir string, perm os.FileMode) error {
	now := time.Now()
	newDir := filepath.Join(dir, strconv.Itoa(now.Year()),
		strconv.Itoa(int(now.Month())))
	return MkdirWhenNo(newDir, perm)
}

// MkdirYMD 在dir下面，按年月日创建文件夹
func MkdirYMD(dir string, perm os.FileMode) error {
	now := time.Now()
	newDir := filepath.Join(dir, strconv.Itoa(now.Year()),
		strconv.Itoa(int(now.Month())),
		strconv.Itoa(now.Day()))
	return MkdirWhenNo(newDir, perm)
}

func DeleteFiles(filePath ...string) error {
	delFailed := make([]string, 0)
	for _, s := range filePath {
		if !Exist(s) {
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
		if !Exist(f) {
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
