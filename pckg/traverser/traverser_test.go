package traverser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TraverserProperlyWorks(t *testing.T) {
	testingDir, _ := os.MkdirTemp("", "testing")

	file1, _ := os.CreateTemp(testingDir, "file.txt")

	testingSubDir, _ := os.MkdirTemp(testingDir, "subdir")

	file2, _ := os.CreateTemp(testingSubDir, "other_file.txt")

	defer os.RemoveAll(testingDir)

	filesChan := GetAllFilesInDirectory(TraversalRequest{TraversalDirectory: testingDir, ExcludedPaths: &[]string{}})
	var resultPath []string
	for path := range filesChan {
		resultPath = append(resultPath, path)
	}

	assert.Equal(t, 2, len(resultPath))
	assert.Equal(t, file1.Name(), resultPath[0])
	assert.Equal(t, file2.Name(), resultPath[1])
}

func Test_IsValidPath(t *testing.T) {
	assert.True(t, isValidPath("/", []string{}))
}

func Test_IsValidPathProperlyUsesIgnores(t *testing.T) {
	assert.False(t, isValidPath("/some/string.txt", []string{"*"}))
	assert.True(t, isValidPath("/some/string.txt", []string{"*.exe"}))
	assert.False(t, isValidPath("/some/string.txt", []string{"/some/*.txt"}))
	assert.True(t, isValidPath("/some/string.txt", []string{"/some/*.md"}))
}

func Test_IsValidPathWithMultipleIgnoreLines(t *testing.T) {
	assert.False(t, isValidPath("/some/string.txt", []string{"*.exe", "*.txt"}))
}

func Test_TraverserWalProperlyIgnoresFiles(t *testing.T) {
	testingDir, _ := os.MkdirTemp("", "testing")

	file1, _ := os.CreateTemp(testingDir, "file.txt")

	testingSubDir, _ := os.MkdirTemp(testingDir, "subdir")

	os.CreateTemp(testingSubDir, "other_file.md")

	defer os.RemoveAll(testingDir)

	filesChan := GetAllFilesInDirectory(TraversalRequest{
		TraversalDirectory: testingDir,
		ExcludedPaths: &[]string{
			testingSubDir,
		},
	})
	var resultPath []string
	for path := range filesChan {
		resultPath = append(resultPath, path)
	}

	assert.Equal(t, 1, len(resultPath))
	assert.Equal(t, file1.Name(), resultPath[0])
}

func Test_TraverserIgnoresDirectiores(t *testing.T) {
	testingDir, _ := os.MkdirTemp("", "testing")

	file1, _ := os.CreateTemp(testingDir, "file.txt")

	testingSubDir, _ := os.MkdirTemp(testingDir, "subdir")
	os.MkdirTemp(testingDir, "otherSubDir")

	os.CreateTemp(testingSubDir, "other_file.md")

	defer os.RemoveAll(testingDir)

	filesChan := GetAllFilesInDirectory(TraversalRequest{
		TraversalDirectory: testingDir,
		ExcludedPaths: &[]string{
			testingSubDir,
		},
	})
	var resultPath []string
	for path := range filesChan {
		resultPath = append(resultPath, path)
	}

	assert.Equal(t, 1, len(resultPath))
	assert.Equal(t, file1.Name(), resultPath[0])
}
