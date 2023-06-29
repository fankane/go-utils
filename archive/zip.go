package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// UnZip 将zipFile 解压缩到 dst 目录, 返回文件名列表
func UnZip(zipFle, dst string) ([]string, error) {
	zr, err := zip.OpenReader(zipFle)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	// 如果解压后不是放在当前目录就按照保存目录去创建目录
	if dst != "" {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return nil, err
		}
	}

	fileNameList := make([]string, 0)
	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		filePath, err := writeFileInZip(dst, file)
		if err != nil {
			return nil, err
		}
		if filePath != "" {
			fileNameList = append(fileNameList, filePath)
		}
	}
	return fileNameList, nil
}

func writeFileInZip(dst string, file *zip.File) (string, error) {
	fileName := file.Name
	if file.Flags == 0 { //处理文件名乱码问题
		//如果标致位是0  则是默认的本地编码   默认为gbk
		i := bytes.NewReader([]byte(file.Name))
		decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
		content, err := ioutil.ReadAll(decoder)
		if err == nil {
			fileName = string(content)
		}
	}

	path := filepath.Join(dst, fileName)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(path, file.Mode()); err != nil {
			return "", err
		}
		return "", nil
	}

	fr, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fr.Close()

	fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
	if err != nil {
		return "", err
	}
	defer fw.Close()

	_, err = io.Copy(fw, fr)
	if err != nil {
		return "", err
	}
	return path, nil
}
