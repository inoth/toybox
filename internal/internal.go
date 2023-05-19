package internal

import (
	"os"
	"path/filepath"
)

func WalkPath(path string) ([]string, error) {
	var files []string
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return files, err
	}
	return files, nil
}
