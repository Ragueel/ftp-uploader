package main

import (
	"fmt"
	"ftp-uploader/cmd/ftp-uploader/cli"

	"github.com/spf13/cobra"
)

const welcomeMessage = `
 $$$$$$\    $$\                                           $$\                           $$\                     
$$  __$$\   $$ |                                          $$ |                          $$ |                    
$$ /  \__|$$$$$$\    $$$$$$\          $$\   $$\  $$$$$$\  $$ | $$$$$$\   $$$$$$\   $$$$$$$ | $$$$$$\   $$$$$$\  
$$$$\     \_$$  _|  $$  __$$\ $$$$$$\ $$ |  $$ |$$  __$$\ $$ |$$  __$$\  \____$$\ $$  __$$ |$$  __$$\ $$  __$$\ 
$$  _|      $$ |    $$ /  $$ |\______|$$ |  $$ |$$ /  $$ |$$ |$$ /  $$ | $$$$$$$ |$$ /  $$ |$$$$$$$$ |$$ |  \__|
$$ |        $$ |$$\ $$ |  $$ |        $$ |  $$ |$$ |  $$ |$$ |$$ |  $$ |$$  __$$ |$$ |  $$ |$$   ____|$$ |      
$$ |        \$$$$  |$$$$$$$  |        \$$$$$$  |$$$$$$$  |$$ |\$$$$$$  |\$$$$$$$ |\$$$$$$$ |\$$$$$$$\ $$ |      
\__|         \____/ $$  ____/          \______/ $$  ____/ \__| \______/  \_______| \_______| \_______|\__|      
                    $$ |                        $$ |                                                            
                    $$ |                        $$ |                                                            
                    \__|                        \__|`

var rootCmd = &cobra.Command{
	Use:   "ftp-uploader",
	Short: "ftp-uploader - a utility to upload files to your ftp server with ignores",
	Long:  "ftp-uploader allows you to upload your files to your ftp server with gitignore like rules",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(welcomeMessage)
		fmt.Println("\n\nUse `ftp-uploader init` to initialize project. Use `ftp-uploader upload` command to upload to ftp.")
	},
}

func init() {
	rootCmd.AddCommand(cli.InitCommand)
	rootCmd.AddCommand(cli.UploadCommand)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
