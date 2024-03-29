package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const uploaderFixture = `
configs:
  default:
    root: .
    uploadRoot: heavy-test
    name: default
    ignorePaths:
      - ftp-uploader.yaml
      - .git/
      - .idea/
`

const uploaderFixtureWithIgnoreFile = `
configs:
  default:
    root: .
    uploadRoot: heavy-test
    name: default
    ignoreFile: %s
    ignorePaths:
      - ftp-uploader.yaml
      - .git/
      - .idea/
`

const mockIgnoreFixture = `
.sample
*.txt
`

func Test_LoadFromFileWorks(t *testing.T) {
	file, err := os.CreateTemp("", "mock.yaml")
	assert.NoError(t, err)
	_, err = file.Write([]byte(uploaderFixture))
	assert.NoError(t, err)

	resultFromFile, err := NewRootFromFile(file.Name(), AuthCredentials{}, 1)
	assert.NoError(t, err)

	settings, ok := resultFromFile.Configs["default"]
	assert.True(t, ok)

	assert.Equal(t, settings.Name, "default")
	assert.Equal(t, settings.UploadRootPath, "heavy-test")
	assert.Equal(t, settings.LocalRootPath, ".")
	assert.Equal(t, len(settings.IgnorePaths), 3)
}

func Test_LoadingFixtureWithGitignoreFile(t *testing.T) {
	ignoreFile, err := os.CreateTemp("", ".ignore")
	assert.NoError(t, err)
	_, err = ignoreFile.Write([]byte(mockIgnoreFixture))
	assert.NoError(t, err)

	file, err := os.CreateTemp("", "mock.yaml")
	assert.NoError(t, err)
	resultFile := fmt.Sprintf(uploaderFixtureWithIgnoreFile, ignoreFile.Name())

	_, err = file.Write([]byte(resultFile))
	assert.NoError(t, err)

	resultFromFile, err := NewRootFromFile(file.Name(), AuthCredentials{}, 1)
	assert.NoError(t, err)

	settings, ok := resultFromFile.Configs["default"]
	assert.True(t, ok)
	assert.Equal(t, 5, len(settings.IgnorePaths))
	assert.Equal(t, ".sample", settings.IgnorePaths[3])
	assert.Equal(t, "*.txt", settings.IgnorePaths[4])
}

func Test_CreateConfigAtPathWorks(t *testing.T) {
	file, err := os.CreateTemp("", "ftp-uploader.yaml")
	assert.NoError(t, err)

	err = CreateDefaultRootFile(file.Name())
	assert.NoError(t, err)
}
