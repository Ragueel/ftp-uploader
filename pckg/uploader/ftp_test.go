package uploader

import (
	"context"
	"ftp-uploader/pckg/config"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var authConfig = config.AuthCredentials{
	Username: "user",
	Password: "password",
	Host:     "localhost:20021",
}

func Test_ProperlyInitializesConnection(t *testing.T) {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ftpUploader, err := NewFtpUploader(timeoutCtx, authConfig, 1)

	assert.NoError(t, err)
	assert.NotNil(t, ftpUploader)

	err = ftpUploader.Close()
	assert.NoError(t, err)
}

func Test_UploadFileAtPathWorks(t *testing.T) {
	todoCtx := context.TODO()
	timeoutCtx, cancel := context.WithTimeout(todoCtx, 2*time.Second)
	defer cancel()

	ftpUploader, _ := NewFtpUploader(todoCtx, authConfig, 1)
	f, err := os.CreateTemp("", "sample.txt")
	f.Write([]byte("Hello world"))

	assert.NoError(t, err)

	defer os.Remove(f.Name())

	task := ftpUploader.UploadFile(timeoutCtx, f.Name(), "test.txt")

	assert.NotNil(t, task)

	for progress := range task.Progress {
		assert.True(t, progress >= 0)
	}

	assert.NoError(t, task.Err)

	conn, err := ftpUploader.getConn(timeoutCtx)
	assert.NoError(t, err)

	file, err := conn.Retr("test.txt")
	assert.NoError(t, err)

	buf, err := io.ReadAll(file)
	assert.NoError(t, err)

	assert.Equal(t, "Hello world", string(buf))
}

func Test_UploadInSubdirectoryWorks(t *testing.T) {
	todoCtx := context.TODO()
	timeoutCtx, cancel := context.WithTimeout(todoCtx, 2*time.Second)
	defer cancel()

	ftpUploader, _ := NewFtpUploader(context.TODO(), authConfig, 1)
	uploadPath := "subdir_sample/test_1/asdasd/test.txt"

	f, err := os.CreateTemp("", "sample.txt")
	f.Write([]byte("Hello world"))

	assert.NoError(t, err)

	defer os.Remove(f.Name())

	task := ftpUploader.UploadFile(timeoutCtx, f.Name(), uploadPath)

	assert.NotNil(t, task)

	for progress := range task.Progress {
		assert.True(t, progress >= 0)
	}

	assert.NoError(t, task.Err)
	conn, err := ftpUploader.getConn(timeoutCtx)
	assert.NoError(t, err)

	file, err := conn.Retr(uploadPath)
	assert.NoError(t, err)

	buf, err := io.ReadAll(file)
	assert.NoError(t, err)

	assert.Equal(t, "Hello world", string(buf))
}
