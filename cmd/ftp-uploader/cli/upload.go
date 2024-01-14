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
	configName      string
	username        string
	password        string
	host            string
	rootConfigPath  string
	connectionCount int
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
	_ = viperEnvs.BindEnv("USERNAME")
	_ = viperEnvs.BindEnv("PASSWORD")
	_ = viperEnvs.BindEnv("HOST")
	_ = viperEnvs.BindEnv("ROOT_CONFIG_PATH")
	_ = viperEnvs.BindEnv("CONNECTION_COUNT")

	viperEnvs.SetDefault("ROOT_CONFIG_PATH", config.DefaultFileName)
	viperEnvs.SetDefault("CONNECTION_COUNT", 1)

	UploadCommand.Flags().StringVarP(&configName, "config", "c", "", "Name of your config")
	UploadCommand.Flags().StringVarP(&username, "username", "u", viperEnvs.GetString("USERNAME"), "Username to use")
	UploadCommand.Flags().StringVarP(&password, "password", "p", viperEnvs.GetString("PASSWORD"), "Password to use")
	UploadCommand.Flags().StringVarP(&host, "host", "s", viperEnvs.GetString("HOST"), "Host server to use")
	UploadCommand.Flags().StringVarP(&rootConfigPath, "root-config-path", "r", viperEnvs.GetString("ROOT_CONFIG_PATH"), "Path to your root config")
	UploadCommand.Flags().IntVarP(&connectionCount, "connection-count", "t", viperEnvs.GetInt("CONNECTION_COUNT"), "Number of parallel connections")
}

func runUpload(_ *cobra.Command, args []string) {
	ctx := context.Background()

	uploadCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(args) > 0 {
		configName = args[0]
	}

	startUploading(uploadCtx, configName)
}

func startUploading(uploadCtx context.Context, uploadingConfig string) {
	fallbackAuthConfig := config.AuthCredentials{
		Username: username,
		Password: password,
		Host:     host,
	}
	rootConfig, err := config.NewRootFromFile(rootConfigPath, fallbackAuthConfig, connectionCount)
	if err != nil {
		fmt.Printf("Invalid root config: %s\n", err)
		return
	}

	if uploadingConfig == "" {
		uploadEveryConfig(uploadCtx, rootConfig)
		return
	}

	uploadConfig, ok := rootConfig.Configs[uploadingConfig]
	if !ok {
		fmt.Printf("Given config is not found: %s\n", uploadingConfig)
		return
	}

	uploadWithConfig(uploadCtx, uploadConfig)
}

func uploadEveryConfig(ctx context.Context, rootConfig *config.Root) {
	for key, val := range rootConfig.Configs {
		fmt.Printf("Uploading: %s\n", key)
		uploadWithConfig(ctx, val)
	}
}

func uploadWithConfig(uploadCtx context.Context, uploadConfig config.UploadSettings) {
	uploadController, err := uploadcontroller.NewFtpUploadController(uploadCtx, *uploadConfig.AuthCredentials, uploadConfig.ConnectionCount)
	if err != nil {
		fmt.Printf("Failed to created uploader: %s\n", err)
		return
	}

	uploadController.UploadFromConfig(uploadCtx, uploadConfig)
}
