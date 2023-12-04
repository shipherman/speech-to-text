package cmd

import (
	"context"
	"fmt"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/internal/clients"
)

type TranscribeServer struct {
	sttservice.UnimplementedSttServiceServer
}

func (t *TranscribeServer) TranscribeAudio(
	audio *sttservice.Audio,
	stream sttservice.SttService_TranscribeAudioServer,
) error {
	text, err := clients.ReqSTT(audio.Audio)
	if err != nil {
		return err
	}

	// TODO
	// Implement pipline for processing audio
	// 1. Check file, if it's a valid wav file
	//		send status ACCEPTED
	// 2. Put file to the worker pool for saving
	// 		send status ORDERED
	// 3. Transcribe file
	// 		send status DONE
	// *send status INVALID on invalid files
	//
	// Additional
	//	convert *.ogg to *.wav [for telegram bot]

	response := &sttservice.Status{Status: sttservice.EnumStatus_STATUS_ACCEPTED}
	stream.Send(response)
	response = &sttservice.Status{Status: sttservice.EnumStatus_STATUS_ORDERED}
	stream.Send(response)
	response = &sttservice.Status{Status: sttservice.EnumStatus_STATUS_DONE}
	stream.Send(response)

	fmt.Println(text)
	return nil
}

func (t *TranscribeServer) GetHistory(ctx context.Context, in *sttservice.User) (*sttservice.History, error) {

	return nil, nil
}
