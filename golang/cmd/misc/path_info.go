package misc

import (
	"os"
	"path/filepath"
	"strings"
)

func pathInfo(path string, extensions []string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileExt := strings.ToLower(filepath.Ext(filePath))
			for _, ext := range extensions {
				if fileExt == ext {
					files = append(files, filePath)
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
