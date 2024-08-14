package source

import (
	"os"
	"path/filepath"
)

func GetFileList(root string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			relativePath = filepath.ToSlash(relativePath) // 将反斜杠替换为正斜杠
			fileList = append(fileList, relativePath)
		}
		return nil
	})

	return fileList, err
}
