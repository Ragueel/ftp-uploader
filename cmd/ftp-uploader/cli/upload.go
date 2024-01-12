package cli

import (
	"context"
	"fmt"
	"ftp-uploader/pckg/config"
	"ftp-uploader/pckg/uploadcontroller"

	"github.com/spf13/cobra"
)

var configName string

var UploadCommand = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"u"},
	Short:   "Upload your config",
	Args:    cobra.RangeArgs(0, 1),
	Run:     runUpload,
}

func init() {
	UploadCommand.Flags().StringVarP(&configName, "config", "c", "", "Name of your config")
}

func runUpload(_ *cobra.Command, args []string) {
	rootConfig, err := config.NewRootConfigFromConfigFile(config.ConfigFileName)
	if err != nil {
		fmt.Printf("Invalid root config: %s\n", err)
		return
	}
	if len(args) > 0 {
		configName = args[0]
	}

	if configName == "" {
		return
	}

	ctx := context.Background()

	uploadCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	uploadConfig, ok := rootConfig.Configs[configName]
	if !ok {
		fmt.Printf("Given config is not found: %s\n", configName)
		return
	}

	uploadWithConfig(uploadCtx, uploadConfig)
}

func uploadWithConfig(uploadCtx context.Context, uploadConfig config.UploadConfig) {
	uploadController, err := uploadcontroller.NewFtpUploadController(uploadCtx, *uploadConfig.AuthConfig)
	if err != nil {
		fmt.Printf("Failed to created uploader: %s", err)
		return
	}

	uploadController.UploadFromConfig(uploadCtx, uploadConfig)
}
