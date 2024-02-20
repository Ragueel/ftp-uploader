package uploader

import (
	"errors"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

var (
	ErrFailedToCreateFile    = errors.New("failed to create file")
	ErrFailedToWriteIntoFile = errors.New("failed to write into file")
)

type SftpUploader struct {
	client *sftp.Client
}

func NewSftp(sshClient *ssh.Client) (*SftpUploader, error) {
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	return &SftpUploader{
		client: client,
	}, nil
}

func (uploader *SftpUploader) UploadFile(ctx context.Context, filePath string, uploadFilePath string) *UploadTask {
	progressChan := make(chan int)
	task := UploadTask{
		Progress: progressChan,
		Err:      nil,
	}
	go func() {
		defer close(progressChan)

		f, err := uploader.client.Create(uploadFilePath)
		if err != nil {
			task.Err = errors.Join(ErrFailedToCreateFile, err)
			return
		}
		defer f.Close()
		_, err = f.Write([]byte("Hello world"))
		if err != nil {
			task.Err = errors.Join(ErrFailedToWriteIntoFile, err)
			return
		}
	}()

	return &task
}
