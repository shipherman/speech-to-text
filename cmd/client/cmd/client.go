package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/pkg/audioconvertor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Send request to server
// Recieve statuses till Done
func SendRequest() error {
	var audio sttservice.Audio
	conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	// sample audio load
	// and check if it has appropriate headers
	audio.Audio = readAudioFromFile()
	_, err = audioconvertor.CheckWAVHeader(audio.Audio)
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
