package cli

import (
	"context"
	"errors"
	"fmt"
	"ftp-uploader/pkg/config"
	"ftp-uploader/pkg/uploadcontroller"
	"os"

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
	Long:    "It is possible to control upload with environement variables. Each variable has a prefix FTP_UPLOADER",
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

	err := startUploading(uploadCtx, configName)
	if err != nil {
		fmt.Printf("Upload failed: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Done")
	}
}

func startUploading(uploadCtx context.Context, uploadingConfig string) error {
	fallbackAuthConfig := config.AuthCredentials{
		Username: username,
		Password: password,
		Host:     host,
	}
	rootConfig, err := config.NewRootFromFile(rootConfigPath, fallbackAuthConfig, connectionCount)
	if err != nil {
		return fmt.Errorf("invalid root config: %w", err)
	}

	if uploadingConfig == "" {
		return uploadEveryConfig(uploadCtx, rootConfig)
	}

	uploadConfig, ok := rootConfig.Configs[uploadingConfig]
	if !ok {
		return errors.New("not found config")
	}

	return uploadWithConfig(uploadCtx, uploadConfig)
}

func uploadEveryConfig(ctx context.Context, rootConfig *config.Root) error {
	for key, val := range rootConfig.Configs {
		fmt.Printf("Uploading: %s\n", key)
		err := uploadWithConfig(ctx, val)
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadWithConfig(uploadCtx context.Context, uploadConfig config.UploadSettings) error {
	uploadController, err := uploadcontroller.NewFtpController(uploadCtx, *uploadConfig.AuthCredentials, uploadConfig.ConnectionCount)
	if err != nil {
		return fmt.Errorf("failed to create uploader: %w", err)
	}

	return uploadController.UploadFromConfig(uploadCtx, uploadConfig)
}
