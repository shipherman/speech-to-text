package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/pkg/audioconverter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client interface {
	SendRequest(context.Context) (string, error)
	Register(ctx context.Context, user string, email string, password string) (int64, error)
	Login(context.Context, string, string) (string, error)
	GetHistory(context.Context) error
}

type STTClient struct {
	sttservice.SttServiceClient
}

var audio sttservice.Audio

func NewClient() (Client, error) {
	var client STTClient
	creds, err := credentials.NewClientTLSFromFile(cfg.CACert, "")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	conn, err := grpc.DialContext(context.Background(), cfg.ServerAddress,
		grpc.WithPerRPCCredentials(fetchToken()),
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return nil, err
	}

	client.SttServiceClient = sttservice.NewSttServiceClient(conn)

	return &client, err
}

// Send request to server
// Recieve statuses till Done
func (c *STTClient) SendRequest(ctx context.Context) (string, error) {
	// perRPC := fetchToken()

	// sample audio load
	// and check if it has appropriate headers
	audio.Audio = readAudioFromFile()
	_, err := audioconverter.CheckWAVHeader(audio.Audio)
	if err != nil {
		return "", err
	}

	stream, err := c.TranscribeAudio(ctx, &audio)
	if err != nil {
		return "", err
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			return "", err
		}
		switch res.Status {
		case sttservice.EnumStatus_STATUS_ACCEPTED:
			log.Println("accepted")
		case sttservice.EnumStatus_STATUS_ORDERED:
			log.Println("orderd")
		case sttservice.EnumStatus_STATUS_PROCESSING:
			log.Println("processing")
		case sttservice.EnumStatus_STATUS_DONE:
			log.Println("DONE")
			return res.Text.Text, nil
		case sttservice.EnumStatus_STATUS_DECLINED:
			return "", fmt.Errorf("audio file could not be processed")
		}
	}
}

// Register registers new user with provided
// username, email and password.
// And returns user ID
func (c *STTClient) Register(ctx context.Context, u, e, p string) (int64, error) {
	respReg, err := c.SttServiceClient.Register(ctx, &sttservice.RegisterRequest{
		Username: u,
		Email:    e,
		Password: p})
	if err != nil {
		return 0, err
	}

	return respReg.UserId, err
}

// Login requests JWT token for provided email and password
func (c *STTClient) Login(ctx context.Context, e, p string) (string, error) {
	fmt.Println("Logiing in")
	respLog, err := c.SttServiceClient.Login(ctx, &sttservice.LoginRequest{
		Email:    e,
		Password: p})
	if err != nil {
		return "", err
	}
	return respLog.Token, nil
}

func (c *STTClient) GetHistory(ctx context.Context) error {
	history, err := c.SttServiceClient.GetHistory(ctx, &sttservice.User{Email: "e"})
	if err != nil {
		return err
	}

	text, err := history.Recv()
	if err != nil {
		return err
	}

	fmt.Println(text)
	return nil
}

// read audio file
func readAudioFromFile() []byte {
	audioBytes, err := os.ReadFile(cfg.FilePath)
	if err != nil {
		fmt.Println(err)
	}
	return audioBytes
}

func fetchToken() *tokenAuth {
	return &tokenAuth{
		token: cfg.AuthToken,
	}
}
