package uploadcontroller

import (
	"context"
	"ftp-uploader/pckg/config"
	"ftp-uploader/pckg/uploader"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockUploader struct {
	filePaths       []string
	uploadFilePaths []string
}

func (u *MockUploader) UploadFile(_ context.Context, filePath string, uploadFilePath string) *uploader.UploadTask {
	progressChan := make(chan int)

	u.filePaths = append(u.filePaths, filePath)
	u.uploadFilePaths = append(u.uploadFilePaths, uploadFilePath)

	task := uploader.UploadTask{
		Progress: progressChan,
		Err:      nil,
	}
	go func() {
		defer close(progressChan)

		progressChan <- 0
		progressChan <- 100
	}()

	return &task
}

func Test_ControllerProperlyUploadsEverything(t *testing.T) {
	testingDir, err := os.MkdirTemp("", "testing")
	assert.NoError(t, err)

	file1, err := os.CreateTemp(testingDir, "test_1.txt")
	assert.NoError(t, err)

	file2, err := os.CreateTemp(testingDir, "test_2.txt")
	assert.NoError(t, err)
	defer os.RemoveAll(testingDir)

	mockUploader := MockUploader{
		filePaths:       []string{},
		uploadFilePaths: []string{},
	}
	uploadController := Controller{Uploader: &mockUploader, maxWorkerCount: 1}

	ctx := context.TODO()
	testingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err = uploadController.UploadFromConfig(testingCtx, config.UploadSettings{
		AuthCredentials: &config.AuthCredentials{},
		LocalRootPath:   testingDir,
		UploadRootPath:  "sample/",
		IgnorePaths:     []string{},
	})

	assert.NoError(t, err)

	assert.Equal(t, 2, len(mockUploader.filePaths))
	assert.Equal(t, 2, len(mockUploader.uploadFilePaths))

	assert.Equal(t, file1.Name(), mockUploader.filePaths[0])
	assert.Equal(t, file2.Name(), mockUploader.filePaths[1])
}
