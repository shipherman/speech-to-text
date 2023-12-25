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

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("register called")
		c, err := NewClient()
		if err != nil {
			log.Println(err)
		}

		t := time.Now().Add(time.Minute * cfg.Timeout)
		ctx, cancel := context.WithDeadline(context.Background(), t)
		defer cancel()
		userID, err := c.Register(ctx, cfg.User, cfg.Email, cfg.Password)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(userID)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	registerCmd.Flags().StringVar(
		&cfg.User,
		"username",
		"user",
		"Username")
	registerCmd.Flags().StringVar(
		&cfg.Email,
		"email",
		"user@mail.ru",
		"Email")
	registerCmd.Flags().StringVar(
		&cfg.Password,
		"password",
		"password",
		"Password")
}
