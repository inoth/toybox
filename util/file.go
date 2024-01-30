package util

import (
	"fmt"
	"os"
	"path/filepath"
	// "github.com/hpcloud/tail"
)

func ReadFile(path string) ([]byte, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents: %s", err.Error())
	}
	return buf, err
}

func WriteFile(path string, buf []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(buf)
	if err != nil {
		return fmt.Errorf("unable to write file: %v", err)
	}
	return nil
}

func WalkPath(path string, wildcards ...string) ([]string, error) {
	wildcard := First("", wildcards)

	var files []string
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if wildcard == "" {
				files = append(files, path)
			} else {
				if matched, err := filepath.Match(wildcard, filepath.Base(path)); matched && err == nil {
					files = append(files, path)
				}
			}
		}
		return nil
	}); err != nil {
		return files, err
	}
	return files, nil
}

func PathGlobPattern(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// func TailFile(path string) {
// 	tail.TailFile(path, tail.Config{Follow: true, ReOpen: true})
// }
