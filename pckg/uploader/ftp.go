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
	PreCreatedDirectories map[string]bool
	authConfig            config.AuthCredentials
	connQueue             chan *ftp.ServerConn
	allConnections        []*ftp.ServerConn
}

func NewFtpUploader(ctx context.Context, authConfig config.AuthCredentials, connectionCount int) (*FtpUploader, error) {
	uploader := FtpUploader{
		allConnections:        make([]*ftp.ServerConn, 0),
		connQueue:             make(chan *ftp.ServerConn, connectionCount),
		PreCreatedDirectories: make(map[string]bool),
	}

	for i := 0; i < connectionCount; i++ {
		connection, err := createConnection(ctx, authConfig)
		if err != nil {
			return nil, err
		}
		uploader.allConnections = append(uploader.allConnections, connection)
		uploader.connQueue <- connection
	}

	return &uploader, nil
}

func createConnection(ctx context.Context, authConfig config.AuthCredentials) (*ftp.ServerConn, error) {
	ftpClient, err := ftp.Dial(authConfig.Host, ftp.DialWithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to dial to remote host: %w", err)
	}
	err = ftpClient.Login(authConfig.Username, authConfig.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return ftpClient, nil
}

func (uploader *FtpUploader) getConn(ctx context.Context) (*ftp.ServerConn, error) {
	select {
	case conn, ok := <-uploader.connQueue:
		if !ok {
			return nil, fmt.Errorf("connection pool closed")
		}
		return conn, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (uploader *FtpUploader) putConn(ctx context.Context, conn *ftp.ServerConn) {
	select {
	case uploader.connQueue <- conn:
		return
	case <-ctx.Done():
		panic(fmt.Errorf("failed to write into connection queue: %w", ctx.Err()))
	}
}

func (uploader *FtpUploader) Close() error {
	close(uploader.connQueue)
	for _, conn := range uploader.allConnections {
		err := conn.Quit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (uploader *FtpUploader) createDirectoryIfNotExists(conn *ftp.ServerConn, uploadPath string) error {
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

		err := conn.MakeDir(remoteDir)
		if err != nil {
			currentDir, err := conn.CurrentDir()
			if err != nil {
				return fmt.Errorf("could not get current directory: %w", err)
			}

			err = conn.ChangeDir(remoteDir)

			if err != nil {
				return fmt.Errorf("failed to create ftp directory: %w", err)
			}

			err = conn.ChangeDir(currentDir)
			if err != nil {
				return fmt.Errorf("failed to reset directory: %w", err)
			}
		}

		uploader.PreCreatedDirectories[remoteDir] = true

	}
	return nil
}

// TODO: Add proper context cancelation
func (uploader *FtpUploader) UploadFile(ctx context.Context, filePath string, uploadFilePath string) *UploadTask {
	progressChan := make(chan int)

	task := UploadTask{
		Progress: progressChan,
		Err:      nil,
	}

	go func() {
		defer close(progressChan)
		ftpConn, err := uploader.getConn(ctx)
		if err != nil {
			task.Err = fmt.Errorf("failed to get connection: %w", err)
			return
		}

		defer uploader.putConn(ctx, ftpConn)

		err = uploader.createDirectoryIfNotExists(ftpConn, uploadFilePath)
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

		err = ftpConn.Stor(uploadFilePath, reader)
		if err != nil {
			task.Err = fmt.Errorf("failed to execute stor on file: %s %w", uploadFilePath, err)
			return
		}
		progressChan <- 100
	}()

	return &task
}
