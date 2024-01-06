package traverser

import (
	"fmt"
	"os"
	"path/filepath"
)

func isValidPath(filePath string, excludedPaths *[]string) bool {
	return true
}

func GetAllFilesInDirectory(traversalInfo TraversalInfo) <-chan string {
	result := make(chan string)

	go func() {
		defer close(result)
		err := filepath.Walk(traversalInfo.TraversalDirectory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !isValidPath(path, traversalInfo.ExcludedPaths) {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			result <- path

			return nil
		})
		if err != nil {
			fmt.Println("Failed to walk directory")
		}
	}()

	return result
}
