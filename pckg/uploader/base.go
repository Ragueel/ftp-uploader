package uploader

type Uploader interface {
	UploadFile(filePath string, uploadFilePath string) *UploadTask
}

type UploadTask struct {
	Progress <-chan int
	Err      error
}
