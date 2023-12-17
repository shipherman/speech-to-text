package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/pkg/audioconverter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type tokenAuth struct {
	token string
}

func (a *tokenAuth) GetRequestMetadata(ctx context.Context,
	uri ...string) (map[string]string, error) {
	return map[string]string{"authorization": a.token, "alg": "HS256"}, nil
}

func (a *tokenAuth) RequireTransportSecurity() bool {
	return true
}

// Send request to server
// Recieve statuses till Done
func SendRequest() error {
	var audio sttservice.Audio
	// perRPC := fetchToken()
	creds, err := credentials.NewClientTLSFromFile("./cert/ca_cert.pem", "x.test.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	// opts := []grpc.DialOption{
	// In addition to the following grpc.DialOption, callers may also use
	// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
	// itself.
	// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
	// oauth.TokenSource requires the configuration of transport
	// credentials.
	// }
	conn, err := grpc.DialContext(context.Background(), cfg.ServerAddress,
		grpc.WithPerRPCCredentials(fetchToken()),
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return err
	}

	// sample audio load
	// and check if it has appropriate headers
	audio.Audio = readAudioFromFile()
	_, err = audioconverter.CheckWAVHeader(audio.Audio)
	if err != nil {
		return err
	}

	ctx := context.Background()
	client := sttservice.NewSttServiceClient(conn)

	stream, err := client.TranscribeAudio(ctx, &audio)
	if err != nil {
		return err
	}
	// respReg, err := client.Register(ctx, &sttservice.RegisterRequest{
	// 	Username: "u1",
	// 	Email:    "e1",
	// 	Password: "p1"})
	// fmt.Println(respReg.UserId, err)

	// fmt.Println("Logiing in")
	// respLog, err := client.Login(ctx, &sttservice.LoginRequest{
	// 	Email:    "e1",
	// 	Password: "p1"})
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(respLog.Token, err)
	// TODO
	// Implement switch case structure for all known statuses
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return err
		}

		fmt.Println(res.Status)
	}
}

// read example file
func readAudioFromFile() []byte {
	audioBytes, err := os.ReadFile(cfg.FilePath)
	if err != nil {
		fmt.Println(err)
	}
	return audioBytes
}

func fetchToken() *tokenAuth {
	return &tokenAuth{
		token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImUxIiwiZXhwIjoxNzAyODUzNDgwLCJ1aWQiOjF9.Yabm0qlMDM3rptOR4oTLctKlrMd9fN4YO7qzSdBhzOk",
	}
}
