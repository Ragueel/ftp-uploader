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

var authConfig = config.AppAuthConfig{
	Username: "user",
	Password: "password",
	Host:     "localhost:20021",
}

func Test_ProperlyInitializesConnection(t *testing.T) {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ftpUploader, err := NewFtpUploader(timeoutCtx, authConfig)
	assert.NoError(t, err)
	assert.NotNil(t, ftpUploader)
	err = ftpUploader.Close()

	assert.NoError(t, err)
}

func Test_UploadFileAtPathWorks(t *testing.T) {
	ftpUploader, _ := NewFtpUploader(context.TODO(), authConfig)
	f, err := os.CreateTemp("", "sample.txt")
	f.Write([]byte("Hello world"))

	assert.NoError(t, err)

	defer os.Remove(f.Name())

	task := ftpUploader.UploadFile(f.Name(), "test.txt")

	assert.NotNil(t, task)

	for progress := range task.Progress {
		assert.True(t, progress >= 0)
	}

	assert.NoError(t, task.Err)

	file, err := ftpUploader.Conn.Retr("test.txt")
	assert.NoError(t, err)

	buf, err := io.ReadAll(file)
	assert.NoError(t, err)

	assert.Equal(t, "Hello world", string(buf))
}
