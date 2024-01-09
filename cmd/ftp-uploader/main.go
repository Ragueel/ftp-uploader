package main

import (
	"context"
	"fmt"
	"ftp-uploader/pckg/config"
	"ftp-uploader/pckg/uploadcontroller"
	"time"
)

func main() {
	ctx := context.Background()

	uploadCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	uploadController, err := uploadcontroller.NewFtpUploadController(uploadCtx, config.AppAuthConfig{
		Username: "user",
		Password: "password",
		Host:     "localhost:20021",
	})
	if err != nil {
		fmt.Printf("Failed to created uploader: %s", err)
		return
	}

	uploadController.UploadFromConfig(uploadCtx, config.UploadConfig{
		LocalRootPath:  ".",
		UploadRootPath: "heavy",
		IgnorePaths:    &[]string{".git/"},
	})
}
