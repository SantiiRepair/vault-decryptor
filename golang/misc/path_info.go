package misc

import (
	"os"
	"path/filepath"
	"strings"
)

func PathInfo(path string, extension string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileExt := strings.ToLower(filepath.Ext(filePath))
			if fileExt == extension {
					files = append(files, filePath)
				
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
