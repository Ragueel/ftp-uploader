package uploadcontroller

import (
	"context"
	"fmt"
	"ftp-uploader/pckg/config"
	"ftp-uploader/pckg/traverser"
	"ftp-uploader/pckg/uploader"
	"ftp-uploader/pckg/worker"
	"strings"
	"sync"
)

type UploadController struct {
	Uploader uploader.Uploader
}

func NewFtpUploadController(ctx context.Context, authConfig config.AuthCredentials) (*UploadController, error) {
	ftpUploader, err := uploader.NewFtpUploader(ctx, authConfig)
	if err != nil {
		return nil, err
	}

	uploadController := &UploadController{Uploader: ftpUploader}
	return uploadController, nil
}

func (uploadController *UploadController) uploadFile(filePath, outputPath string) (interface{}, error) {
	fmt.Printf("Uploading file: %s\n", filePath)
	result := uploadController.Uploader.UploadFile(filePath, outputPath)

	for progress := range result.Progress {
		fmt.Printf("Progress: %d\n", progress)
	}
	if result.Err != nil {
		return nil, result.Err
	}

	return filePath, nil
}

func (uploadController *UploadController) UploadFromConfig(ctx context.Context, conf config.UploadSettings) {
	filesChan := traverser.GetAllFilesInDirectory(traverser.TraversalRequest{
		TraversalDirectory: conf.LocalRootPath, ExcludedPaths: conf.IgnorePaths,
	})

	uploadWorker := worker.NewPool(1)

	uploadJobsChan := make(chan worker.Job, uploadWorker.WorkersCount)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(uploadJobsChan)

		for filePath := range filesChan {
			job := worker.Job{
				Descriptor: fmt.Sprintf("FileUploaderJob: %s\n", filePath),
				ExecFn: func(path string) worker.ExecutionFn {
					return func(ctx context.Context) (interface{}, error) {
						trimmedFilePath := strings.TrimPrefix(path, conf.LocalRootPath)
						uploadDestination := fmt.Sprintf("%s/%s", conf.UploadRootPath, trimmedFilePath)

						return uploadController.uploadFile(path, uploadDestination)
					}
				}(filePath),
			}
			uploadJobsChan <- job
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer uploadWorker.Close()

		uploadWorker.Run(ctx, uploadJobsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for result := range uploadWorker.Results {
			if result.Err != nil {
				fmt.Printf("Failed task: %s\n", result.Err)
				continue
			}

			fmt.Printf("Completed task: %v\n", result.Val)
		}
	}()

	wg.Wait()

	fmt.Println("Completed upload of config")
}
