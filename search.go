package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func search(dirname string, filename string) string {
	var result string
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == filename {
			result = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}
	if result == "" {
		return "File not found"
	}
	return result
}

func main() {
	directory := "."
	filename := "test.txt"

	fmt.Print(search(directory, filename))
}
