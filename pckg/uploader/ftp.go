package uploader

import (
	"bufio"
	"context"
	"fmt"
	"ftp-uploader/pckg/config"
	"os"

	"github.com/jlaffaye/ftp"
)

type Uploader interface {
	UploadFile(filePath string, uploadFilePath string) UploadTask
}

type UploadTask struct {
	Progress <-chan int
	Err      error
}

type FtpUploader struct {
	Conn *ftp.ServerConn
}

func NewFtpUploader(ctx context.Context, authConfig config.AppAuthConfig) (*FtpUploader, error) {
	ftpClient, err := ftp.Dial(authConfig.Host, ftp.DialWithContext(ctx))
	if err != nil {
		return nil, err
	}
	err = ftpClient.Login(authConfig.Username, authConfig.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return &FtpUploader{Conn: ftpClient}, nil
}

func (uploader *FtpUploader) UploadFile(filePath string, uploadFilePath string) UploadTask {
	progressChan := make(chan int)

	task := UploadTask{
		Progress: progressChan,
		Err:      nil,
	}

	go func() {
		defer close(progressChan)

		file, err := os.Open(filePath)
		if err != nil {
			task.Err = fmt.Errorf("could not open os file: %w", err)
			return
		}
		reader := bufio.NewReader(file)
		progressChan <- 0

		err = uploader.Conn.Stor(uploadFilePath, reader)
		if err != nil {
			task.Err = fmt.Errorf("failed to execute stor on file: %w", err)
			return
		}
		progressChan <- 100
	}()

	return task
}
