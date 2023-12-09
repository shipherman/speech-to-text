/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/shipherman/speech-to-text/gen/ent"
	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/internal/clients"
	"github.com/shipherman/speech-to-text/internal/db"
	"github.com/shipherman/speech-to-text/internal/transport/routes"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type Config struct {
	ServerAddress string
	DSN           string
	STTAddress    string
}

var cfg Config
var DBConn db.Connector

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "SpeechToText service. Accept wav audio files and returns text.",
	Long:  `SpeechToText service. Accept wav audio files and returns text.`,
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

	// Configure clients to connect to other services
	// STT
	clients.ConfigureSTT("http://localhost:9090", time.Second*5)

	// err = clients.ReqSTT("/home/tas/Downloads/sample.wav")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Server configuration
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: routes.Router,
	}

	// Listener configuration for gRPC connection
	tcpListen, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()

	// Init connection to DB
	// Fatal on error
	client, err := ent.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}

	// Register STT Server with Transcribe server instance
	sttservice.RegisterSttServiceServer(grpcServer,
		&TranscribeServer{DBClient: db.Connector{Client: client}})

	// Run http and grpc server
	for {
		log.Fatal(grpcServer.Serve(tcpListen))
		log.Fatal(server.ListenAndServe())
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.server.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().
		StringVarP(&cfg.ServerAddress,
			"server-address",
			"a",
			"localhost:8282",
			"GRPC server address to run")
	rootCmd.PersistentFlags().
		StringVarP(&cfg.STTAddress,
			"stt-address",
			"s",
			"localhost:9090",
			"Address to connect to Speech-to-text neural network service")
	rootCmd.PersistentFlags().
		StringVarP(&cfg.DSN,
			"dsn",
			"d",
			"host=127.0.0.1 port=5432 user=postgres password=pass dbname=postgres sslmode=disable",
			"Postgres Database connection string")

	// Configure db schema
	err := db.ConfigureSchema(cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
}
