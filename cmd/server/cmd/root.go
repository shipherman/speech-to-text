/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/shipherman/speech-to-text/gen/ent"
	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/internal/clients"
	"github.com/shipherman/speech-to-text/internal/db"
	"github.com/shipherman/speech-to-text/internal/logger"
	"github.com/shipherman/speech-to-text/internal/services/auth"
	"github.com/shipherman/speech-to-text/pkg/fsstore"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	ServerAddress string
	DSN           string
	STTAddress    string
	StorePath     string
	Secret        string
	AuthTimeOut   time.Duration
	ServerCert    string
	ServerKey     string
}

var cfg Config

var grpcServer *grpc.Server
var DBConn db.Connector
var idleConnectionsClosed = make(chan struct{})
var transcribeServer *TranscribeServer

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
	clients.ConfigureSTT(cfg.STTAddress, time.Second*5)

	// Listener configuration for gRPC connection
	tcpListen, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Init connection to DB
	// Fatal on error
	client, err := ent.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	// Init db connector
	dbclient := db.Connector{Client: client}

	// Init blob store
	fsstore := fsstore.NewFSStore(cfg.StorePath)

	// Auth interceptor initiation
	servAuth := auth.New(&dbclient, &dbclient, cfg.AuthTimeOut, cfg.Secret)

	// Load server certificate
	cert, err := tls.LoadX509KeyPair(cfg.ServerCert, cfg.ServerKey)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	// Define gRPC server options
	// Authenticator + ATLS creds
	opts := []grpc.ServerOption{
		// grpc.ChainUnaryInterceptor(auth.AuthUnaryInterceptor),
		grpc.ChainStreamInterceptor(servAuth.AuthStreamInterceptor, grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(logger.ZapInterceptor()))),
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.MaxConcurrentStreams(1),
	}
	grpcServer = grpc.NewServer(opts...)

	transcribeServer = &TranscribeServer{
		DBClient: dbclient,
		auth:     *servAuth,
		Store:    fsstore,
	}
	// Register STT Server with Transcribe server instance
	sttservice.RegisterSttServiceServer(grpcServer,
		transcribeServer)

	// Run http and grpc server
	go gracefullShutdown()

	log.Fatal(grpcServer.Serve(tcpListen))

	<-idleConnectionsClosed

}

func gracefullShutdown() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigint
	log.Println("Shutting down server")

	transcribeServer.DBClient.Close()
	transcribeServer.Store.Close()
	grpcServer.GracefulStop()

	close(idleConnectionsClosed)

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
			"http://localhost:9090",
			"Address to connect to Speech-to-text neural network service")
	rootCmd.PersistentFlags().
		StringVarP(&cfg.DSN,
			"dsn",
			"d",
			"host=127.0.0.1 port=5432 user=postgres password=pass dbname=postgres sslmode=disable",
			"Postgres Database connection string")
	rootCmd.PersistentFlags().
		StringVarP(&cfg.StorePath,
			"store-path",
			"p",
			"/tmp/stt/store",
			"Path to local blob storage")
	rootCmd.PersistentFlags().
		StringVar(&cfg.Secret,
			"secret",
			"verysecretstring",
			"Secret string to generete JWT")
	rootCmd.PersistentFlags().
		StringVar(&cfg.ServerCert,
			"server-certificate",
			"./cert_test/server_cert.pem",
			"Path to server certificate")
	rootCmd.PersistentFlags().
		StringVar(&cfg.ServerKey,
			"server-key",
			"./cert_test/server_key.pem",
			"Path to server key")
	rootCmd.PersistentFlags().
		DurationVar(&cfg.AuthTimeOut,
			"auth-timeout",
			time.Hour*3,
			"Timeout for authenticator server provides to clients")

		// Configure db schema
	err := db.ConfigureSchema(cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
}
