package uploader

import "golang.org/x/net/context"

type Uploader interface {
	UploadFile(ctx context.Context, filePath string, uploadFilePath string) *UploadTask
}

type UploadTask struct {
	Progress <-chan int
	Err      error
}
