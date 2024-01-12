package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const ConfigFileName = "ftp-uploader.yaml"

type AppAuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type UploadConfig struct {
	AuthConfig     *AppAuthConfig `yaml:"authConfig,omitempty"`
	LocalRootPath  string         `yaml:"root"`
	UploadRootPath string         `yaml:"uploadRoot"`
	Name           string         `yaml:"name,omitempty"`
	IgnoreFile     string         `yaml:"ignoreFile,omitempty"`
	IgnorePaths    []string       `yaml:"ignorePaths"`
}

type RootConfig struct {
	Configs map[string]UploadConfig `yaml:"configs"`
}

func NewEmptyUploadConfig() UploadConfig {
	return UploadConfig{
		LocalRootPath:  ".",
		UploadRootPath: "my-relative-path/",
		Name:           "default",
		IgnorePaths:    []string{"ftp-uploader.yaml"},
	}
}

func NewEmptyRootConfig() RootConfig {
	return RootConfig{Configs: map[string]UploadConfig{"default": NewEmptyUploadConfig()}}
}

func NewRootConfigFromConfigFile(configPath string) (*RootConfig, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	rootConfig := RootConfig{}
	err = yaml.Unmarshal(file, &rootConfig)
	if err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}

	for name, uploadConfig := range rootConfig.Configs {
		if uploadConfig.AuthConfig == nil {
			uploadConfig.AuthConfig = NewAuthConfigFromEnv()
		}
		if uploadConfig.IgnoreFile != "" {
			ignoreLines, err := readIgnoreFile(uploadConfig.IgnoreFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read ignore file of %s: %w", name, err)
			}
			uploadConfig.IgnorePaths = append(uploadConfig.IgnorePaths, ignoreLines...)
		}

		uploadConfig.UploadRootPath = strings.TrimSuffix(uploadConfig.UploadRootPath, "/")
		rootConfig.Configs[name] = uploadConfig
	}

	return &rootConfig, nil
}

func readIgnoreFile(ignoreFilePath string) ([]string, error) {
	var ignoreLines []string
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		return ignoreLines, fmt.Errorf("could not read ignore file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ignoreLines = append(ignoreLines, scanner.Text())
	}

	return ignoreLines, nil
}

func NewAuthConfigFromEnv() *AppAuthConfig {
	viper.SetEnvPrefix("FTP_UPLOADER")
	viper.BindEnv("USERNAME")
	viper.BindEnv("PASSWORD")
	viper.BindEnv("HOST")
	return &AppAuthConfig{
		Username: viper.GetString("USERNAME"),
		Password: viper.GetString("PASSWORD"),
		Host:     viper.GetString("HOST"),
	}
}

func NewAuthConfigFromParams() AppAuthConfig {
	return AppAuthConfig{}
}

func NewUploadConfig() UploadConfig {
	return UploadConfig{}
}
