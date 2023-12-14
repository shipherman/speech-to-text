package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/pkg/audioconverter"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

// Send request to server
// Recieve statuses till Done
func SendRequest() error {
	var audio sttservice.Audio
	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(fetchToken())}
	creds, err := credentials.NewClientTLSFromFile("./cert/ca_cert.pem", "x.test.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPC),
		// oauth.TokenSource requires the configuration of transport
		// credentials.
		grpc.WithTransportCredentials(creds),
	}
	conn, err := grpc.Dial(cfg.ServerAddress, opts...)
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

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}
