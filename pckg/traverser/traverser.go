package traverser

import (
	"fmt"
	"os"
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
)

func isValidPath(filePath string, excludedPaths []string) bool {
	if len(excludedPaths) == 0 {
		return true
	}
	ignoreConfig := ignore.CompileIgnoreLines(excludedPaths...)

	return !ignoreConfig.MatchesPath(filePath)
}

func GetAllFilesInDirectory(traversalRequest TraversalRequest) <-chan string {
	result := make(chan string)

	go func() {
		defer close(result)
		err := filepath.Walk(traversalRequest.TraversalDirectory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !isValidPath(path, traversalRequest.ExcludedPaths) {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			result <- path

			return nil
		})
		if err != nil {
			fmt.Println("Failed to walk directory: ", err)
		}
	}()

	return result
}
