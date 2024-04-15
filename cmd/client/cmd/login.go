package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shipherman/speech-to-text/internal/client"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to STT service",
	Long:  `Login command returns JWT token to use it for *transcribe* subcommand`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewClient(cfg)
		if err != nil {
			log.Println(err)
		}

		t := time.Now().Add(time.Minute * cfg.Timeout)
		ctx, cancel := context.WithDeadline(context.Background(), t)
		defer cancel()
		authToken, err := c.Login(ctx, cfg.Email, cfg.Password)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(authToken)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loginCmd.Flags().StringVar(
		&cfg.Email,
		"email",
		"user@mail.ru",
		"Email")
	loginCmd.Flags().StringVar(
		&cfg.Password,
		"password",
		"password",
		"Password")
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
