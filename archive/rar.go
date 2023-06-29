package archive

import (
	"path/filepath"

	"github.com/fankane/go-utils/file"
	"github.com/mholt/archiver"
)

func UnRar(rarFile, dst string) ([]string, error) {
	r := archiver.NewRar()
	if err := r.Unarchive(rarFile, dst); err != nil {
		return nil, err
	}
	fileNames, err := file.DirFiles(dst)
	if err != nil {
		return nil, err
	}
	fileFullPath := make([]string, 0, len(fileNames))
	for _, name := range fileNames {
		fileFullPath = append(fileFullPath, filepath.Join(dst, name))
	}
	return fileFullPath, nil
}
