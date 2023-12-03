package cmd

import (
	"context"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/internal/clients"
)

type TranscribeServer struct {
	sttservice.UnimplementedSttServiceServer
}

func (t *TranscribeServer) TranscribeAudio(*sttservice.Audio, sttservice.SttService_TranscribeAudioServer) error {
	clients.ReqSTT("/tmp/stt/audio.wav")

	return nil
}

func (t *TranscribeServer) GetHistory(ctx context.Context, in *sttservice.User) (*sttservice.History, error) {

	return nil, nil
}
