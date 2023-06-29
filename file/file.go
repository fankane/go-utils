package file

import (
	"os"
)

// DirFiles list of filenames in the directory
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
