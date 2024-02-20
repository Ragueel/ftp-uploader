package cli

import (
	"fmt"
	"ftp-uploader/pkg/config"
	"os"

	"github.com/spf13/cobra"
)

var InitCommand = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Init project",
	Run:     runInit,
}

func runInit(_ *cobra.Command, _ []string) {
	err := config.CreateDefaultRootFile(config.DefaultFileName)
	if err != nil {
		fmt.Printf("Failed initializing file: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Initialized default config at: %s\n", config.DefaultFileName)
	}
}
