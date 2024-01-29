package uploadcontroller

import (
	"context"
	"errors"
	"fmt"
	"ftp-uploader/pckg/config"
	"ftp-uploader/pckg/traverser"
	"ftp-uploader/pckg/uploader"
	"ftp-uploader/pckg/worker"
	"strings"
	"sync"
)

type Controller struct {
	maxWorkerCount int
	Uploader       uploader.Uploader
}

func NewFtpController(ctx context.Context, authConfig config.AuthCredentials, connectionCount int) (*Controller, error) {
	ftpUploader, err := uploader.NewFtpUploader(ctx, authConfig, connectionCount)
	if err != nil {
		return nil, err
	}
	if connectionCount < 1 {
		return nil, fmt.Errorf("invalid number of connections should be more than 0, was given: %d", connectionCount)
	}

	uploadController := &Controller{Uploader: ftpUploader, maxWorkerCount: connectionCount}

	return uploadController, nil
}

// TODO: add retries
func (uploadController *Controller) uploadFile(ctx context.Context, filePath, outputPath string) (interface{}, error) {
	fmt.Printf("Uploading file: %s\n", filePath)
	result := uploadController.Uploader.UploadFile(ctx, filePath, outputPath)

	for progress := range result.Progress {
		fmt.Printf("Progress: %d\n", progress)
	}
	if result.Err != nil {
		return nil, result.Err
	}

	return filePath, nil
}

func (uploadController *Controller) UploadFromConfig(ctx context.Context, conf config.UploadSettings) error {
	if uploadController.maxWorkerCount < 1 {
		return errors.New("invalid number of workers")
	}

	filesChan := traverser.GetAllFilesInDirectory(traverser.TraversalRequest{
		TraversalDirectory: conf.LocalRootPath, ExcludedPaths: conf.IgnorePaths,
	})

	uploadWorker := worker.NewPool(uploadController.maxWorkerCount)

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
						// TODO: add better trim handling in case of name collisions
						if conf.LocalRootPath == "." {
							uploadDestination = fmt.Sprintf("%s/%s", conf.UploadRootPath, path)
						}

						return uploadController.uploadFile(ctx, path, uploadDestination)
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

	return nil
}
