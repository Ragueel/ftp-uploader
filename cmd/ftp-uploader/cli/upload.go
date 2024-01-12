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
	configName     string
	username       string
	password       string
	host           string
	rootConfigPath string
)

var UploadCommand = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"u"},
	Short:   "Upload your config",
	Args:    cobra.RangeArgs(0, 1),
	Run:     runUpload,
}

func init() {
	viperEnvs := viper.New()
	viperEnvs.SetEnvPrefix("FTP_UPLOADER")
	viperEnvs.BindEnv("USERNAME")
	viperEnvs.BindEnv("PASSWORD")
	viperEnvs.BindEnv("HOST")
	viperEnvs.BindEnv("ROOT_CONFIG_PATH")
	viperEnvs.SetDefault("ROOT_CONFIG_PATH", config.ConfigFileName)

	UploadCommand.Flags().StringVarP(&configName, "config", "c", "", "Name of your config")
	UploadCommand.Flags().StringVarP(&username, "username", "u", viperEnvs.GetString("USERNAME"), "Username to use")
	UploadCommand.Flags().StringVarP(&password, "password", "p", viperEnvs.GetString("PASSWORD"), "Password to use")
	UploadCommand.Flags().StringVarP(&host, "host", "s", viperEnvs.GetString("HOST"), "Host server to use")
	UploadCommand.Flags().StringVarP(&rootConfigPath, "root-config-path", "r", viperEnvs.GetString("ROOT_CONFIG_PATH"), "Path to your root config")
}

func runUpload(_ *cobra.Command, args []string) {
	fallbackAuthConfig := config.AuthCredentials{
		Username: username,
		Password: password,
		Host:     host,
	}
	rootConfig, err := config.NewRootFromFile(rootConfigPath, fallbackAuthConfig)
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