package cli

import (
	"fmt"
	"ftp-uploader/pckg/config"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var InitCommand = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Init project",
	Run:     runInit,
}

func runInit(_ *cobra.Command, _ []string) {
	rootConfig := config.NewEmptyRoot()

	result, err := yaml.Marshal(&rootConfig)
	if err != nil {
		fmt.Printf("Failed to initialize empty config project %\n", err)
		return
	}

	file, err := os.Create(config.DefaultFileName)
	if err != nil {
		fmt.Printf("Failed to create wile %s\n", err)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Could not close file")
		}
	}()

	_, err = file.Write([]byte(string(result)))
	if err != nil {
		fmt.Printf("Could not write into file: %s\n", err)
		return
	}

	fmt.Println("Initialized default `ftp-uploader.yaml` file")
}
