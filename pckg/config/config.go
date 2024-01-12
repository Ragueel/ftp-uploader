package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = "ftp-uploader.yaml"

type AuthCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type UploadSettings struct {
	AuthCredentials *AuthCredentials `yaml:"authConfig,omitempty"`
	LocalRootPath   string           `yaml:"root"`
	UploadRootPath  string           `yaml:"uploadRoot"`
	Name            string           `yaml:"name,omitempty"`
	IgnoreFile      string           `yaml:"ignoreFile,omitempty"`
	IgnorePaths     []string         `yaml:"ignorePaths"`
}

type Root struct {
	Configs map[string]UploadSettings `yaml:"configs"`
}

func NewEmptyUploadSettings() UploadSettings {
	return UploadSettings{
		LocalRootPath:  ".",
		UploadRootPath: "my-relative-path/",
		Name:           "default",
		IgnorePaths:    []string{"ftp-uploader.yaml"},
	}
}

func NewEmptyRoot() Root {
	return Root{Configs: map[string]UploadSettings{"default": NewEmptyUploadSettings()}}
}

func NewRootFromFile(configPath string, fallbackAuth AuthCredentials) (*Root, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	rootConfig := Root{}
	err = yaml.Unmarshal(file, &rootConfig)
	if err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}

	for name, uploadConfig := range rootConfig.Configs {
		if uploadConfig.AuthCredentials == nil {
			uploadConfig.AuthCredentials = &fallbackAuth
		}
		if uploadConfig.IgnoreFile != "" {
			ignoreLines, err := readIgnoreFile(uploadConfig.IgnoreFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read ignore file of %s: %w", name, err)
			}
			uploadConfig.IgnorePaths = append(uploadConfig.IgnorePaths, ignoreLines...)
		}

		uploadConfig.Name = name
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
