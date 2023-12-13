/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// TODO
// Add JSON configuration support
type Configuratoin struct {
	ServerAddress string `json:"server_address"`
	FilePath      string
}

var cfg Configuratoin

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "Speech-to-text client",
	Long: `Speech-to-text client allows to transcribe audio files to text
			via remote NN service`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	// sample request
	err = SendRequest()
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.client.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(
		&cfg.ServerAddress,
		"server_address",
		"a",
		"localhost:8282",
		"STT server address")
	rootCmd.PersistentFlags().StringVarP(
		&cfg.FilePath,
		"filepath",
		"f",
		"/tmp/stt/ru-Peacock.wav",
		"Path to file to process")
}
