package uploader

import (
	"bufio"
	"context"
	"fmt"
	"ftp-uploader/pckg/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/jlaffaye/ftp"
)

type FtpUploader struct {
	Conn                  *ftp.ServerConn
	PreCreatedDirectories map[string]bool
}

func NewFtpUploader(ctx context.Context, authConfig config.AppAuthConfig) (*FtpUploader, error) {
	ftpClient, err := ftp.Dial(authConfig.Host, ftp.DialWithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to dial to remote host: %w", err)
	}
	err = ftpClient.Login(authConfig.Username, authConfig.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	// TODO: Add connection pooling
	return &FtpUploader{Conn: ftpClient, PreCreatedDirectories: make(map[string]bool)}, nil
}

func (uploader *FtpUploader) Close() error {
	return uploader.Conn.Quit()
}

func (uploader *FtpUploader) createDirectoryIfNotExists(uploadPath string) error {
	uploadFilePathDir := filepath.Dir(uploadPath)

	if uploadFilePathDir == "." {
		return nil
	}

	directories := strings.Split(uploadFilePathDir, "/")

	for i := range directories {
		remoteDir := filepath.Join(directories[:i+1]...)
		if uploader.PreCreatedDirectories[remoteDir] {
			continue
		}

		err := uploader.Conn.MakeDir(remoteDir)
		if err != nil {
			currentDir, err := uploader.Conn.CurrentDir()
			if err != nil {
				return fmt.Errorf("could not get current directory: %w", err)
			}

			err = uploader.Conn.ChangeDir(remoteDir)

			if err != nil {
				return fmt.Errorf("failed to create ftp directory: %w", err)
			}

			err = uploader.Conn.ChangeDir(currentDir)
			if err != nil {
				return fmt.Errorf("failed to reset directory: %w", err)
			}
		}

		uploader.PreCreatedDirectories[remoteDir] = true

	}
	return nil
}

// TODO: Add proper context cancelation
func (uploader *FtpUploader) UploadFile(filePath string, uploadFilePath string) *UploadTask {
	progressChan := make(chan int)

	task := UploadTask{
		Progress: progressChan,
		Err:      nil,
	}

	go func() {
		defer close(progressChan)
		err := uploader.createDirectoryIfNotExists(uploadFilePath)
		if err != nil {
			task.Err = err
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			task.Err = fmt.Errorf("could not open os file: %w", err)
			return
		}
		reader := bufio.NewReader(file)

		progressChan <- 0

		err = uploader.Conn.Stor(uploadFilePath, reader)
		if err != nil {
			task.Err = fmt.Errorf("failed to execute stor on file: %s %w", uploadFilePath, err)
			return
		}
		progressChan <- 100
	}()

	return &task
}
