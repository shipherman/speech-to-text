/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// historyCmd represents the transcribe command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := NewClient()
		if err != nil {
			log.Println(err)
		}

		t := time.Now().Add(time.Minute * cfg.Timeout)
		ctx, cancel := context.WithDeadline(context.Background(), t)
		defer cancel()
		err = c.GetHistory(ctx)
		if err != nil {
			log.Println(err)
		}

		// Show the result
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// historyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// historyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	historyCmd.Flags().StringVarP(
		&cfg.FilePath,
		"filepath",
		"f",
		"/tmp/stt/ru-Peacock.wav",
		"Path to file to process")
	historyCmd.Flags().StringVar(
		&cfg.AuthToken,
		"authtoken",
		"",
		"Authentication token string")
	historyCmd.MarkFlagRequired("authtoken")
}
