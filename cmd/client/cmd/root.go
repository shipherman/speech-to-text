package cmd

import (
	"os"
	"time"

	"github.com/shipherman/speech-to-text/internal/client"
	"github.com/spf13/cobra"
)

var cfg client.Options

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "Speech-to-text client",
	Long:  `Speech-to-text client allows to transcribe audio files to text via remote NN service`,
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
	rootCmd.PersistentFlags().StringVar(
		&cfg.CACert,
		"cacert",
		"./cert_test/ca_cert.pem",
		"Path to CA certificate")
	rootCmd.PersistentFlags().DurationVar(
		&cfg.Timeout,
		"timeout",
		time.Minute*1,
		"Command execution timeout")
}
