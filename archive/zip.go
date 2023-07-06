package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fankane/go-utils/file"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// CreateZip 创建zip包，zipFileName 是完整路径带文件名, fileList 是要加入的文件列表
func CreateZip(zipFileName string, fileList []string) error {
	if err := checkZipFile(zipFileName); err != nil {
		return err
	}
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("create zip file:[]%s err:%s", zipFileName, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	for _, file := range fileList {
		if err = writeFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func ZipDir(zipFileName, dir string) (err error) {
	if err := checkZipFile(zipFileName); err != nil {
		return err
	}
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("create zip file:[]%s err:%s", zipFileName, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	//因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(dir, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		// 替换文件信息中的文件名
		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zipWriter.CreateHeader(fh)
		if err != nil {
			return
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return
		}

		// 将打开的文件 Copy 到 w
		_, err = io.Copy(w, fr)
		if err != nil {
			return
		}
		return nil
	})
}

func checkZipFile(zipFileName string) error {
	if file.FileExist(zipFileName) {
		return fmt.Errorf("zipFile:%s already exist", zipFileName)
	}
	if !strings.HasSuffix(zipFileName, ".zip") {
		return fmt.Errorf("zipfile:%s must be *.zip format", zipFileName)
	}
	return nil
}

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
		filePath, err := saveFileInZip(dst, file)
		if err != nil {
			return nil, err
		}
		if filePath != "" {
			fileNameList = append(fileNameList, filePath)
		}
	}
	return fileNameList, nil
}

func writeFileToZip(zipWriter *zip.Writer, file string) error {
	fileContent, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileContent.Close()
	zipEntry, err := zipWriter.Create(filepath.Base(file))
	if err != nil {
		return fmt.Errorf("zip writer create failed err:%s", err)
	}
	_, err = io.Copy(zipEntry, fileContent)
	if err != nil {
		return fmt.Errorf("write file content err:%s", err)
	}
	return nil
}

func saveFileInZip(dst string, file *zip.File) (string, error) {
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
