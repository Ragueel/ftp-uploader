/*
Entry point for all commands
*/
package main

import (
	"fmt"
	"ftp-uploader/cmd/ftp-uploader/cli"
	"os"

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
                    \__|                        \__|

Use 'ftp-uploader init' to initialize a project. Use 'ftp-uploader upload' to upload files to your ftp server. 

Read more at: https://github.com/Ragueel/ftp-uploader `

var rootCmd = &cobra.Command{
	Use:   "ftp-uploader",
	Short: "ftp-uploader allows you to upload your files to your ftp server with gitignore like logic",
	Long:  "ftp-uploader allows you to upload your files to your ftp server with gitignore like logic",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(welcomeMessage)
	},
}

func init() {
	rootCmd.AddCommand(cli.InitCommand)
	rootCmd.AddCommand(cli.UploadCommand)
}

func main() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Failed %s\n", err)
		os.Exit(1)
	}
}
