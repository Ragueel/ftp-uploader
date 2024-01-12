package cli

import (
	"context"
	"fmt"
	"ftp-uploader/pckg/config"
	"ftp-uploader/pckg/uploadcontroller"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configName string
	username   string
	password   string
	host       string
)

var UploadCommand = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"u"},
	Short:   "Upload your config",
	Args:    cobra.RangeArgs(0, 1),
	Run:     runUpload,
}

func init() {
	viper.SetEnvPrefix("FTP_UPLOADER")
	viper.BindEnv("USERNAME")
	viper.BindEnv("PASSWORD")
	viper.BindEnv("HOST")

	UploadCommand.Flags().StringVarP(&configName, "config", "c", "", "Name of your config")
	UploadCommand.Flags().StringVarP(&username, "username", "u", viper.GetString("USERNAME"), "Username to use")
	UploadCommand.Flags().StringVarP(&password, "password", "p", viper.GetString("PASSWORD"), "Password to use")
	UploadCommand.Flags().StringVarP(&host, "host", "s", viper.GetString("HOST"), "Host server to use")
}

func runUpload(_ *cobra.Command, args []string) {
	fallbackAuthConfig := config.AuthCredentials{
		Username: username,
		Password: password,
		Host:     host,
	}
	rootConfig, err := config.NewRootFromFile(config.ConfigFileName, fallbackAuthConfig)
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

func uploadWithConfig(uploadCtx context.Context, uploadConfig config.UploadSettings) {
	uploadController, err := uploadcontroller.NewFtpUploadController(uploadCtx, *uploadConfig.AuthCredentials)
	if err != nil {
		fmt.Printf("Failed to created uploader: %s", err)
		return
	}

	uploadController.UploadFromConfig(uploadCtx, uploadConfig)
}
