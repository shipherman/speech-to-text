/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// transcribeCmd represents the transcribe command
var transcribeCmd = &cobra.Command{
	Use:   "transcribe",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("transcribe called")
		c, err := NewClient()
		if err != nil {
			log.Println(err)
		}

		t := time.Now().Add(time.Minute * cfg.Timeout)
		ctx, cancel := context.WithDeadline(context.Background(), t)
		defer cancel()
		c.SendRequest(ctx)
	},
}

func init() {
	rootCmd.AddCommand(transcribeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transcribeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transcribeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	transcribeCmd.Flags().StringVarP(
		&cfg.FilePath,
		"filepath",
		"f",
		"/tmp/stt/ru-Peacock.wav",
		"Path to file to process")
	transcribeCmd.Flags().StringVar(
		&cfg.AuthToken,
		"authtoken",
		"",
		"Authentication token string")
}
