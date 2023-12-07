package cmd

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/internal/clients"
	"github.com/shipherman/speech-to-text/pkg/audioconverter"
	"github.com/shipherman/speech-to-text/pkg/fsstore"
)

type TranscribeServer struct {
	sttservice.UnimplementedSttServiceServer
}

var LocStore *fsstore.FSStore

func (t *TranscribeServer) TranscribeAudio(
	audio *sttservice.Audio,
	stream sttservice.SttService_TranscribeAudioServer,
) error {
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

	var response *sttservice.Status

	// Check wave header and if its valid send status ACCEPTED
	_, err := audioconverter.CheckWAVHeader(audio.Audio)
	if err != nil {
		return err
	}
	response = &sttservice.Status{Status: sttservice.EnumStatus_STATUS_ACCEPTED}
	stream.Send(response)

	// Calculate hash sum
	h := md5.New()
	_, err = h.Write(audio.Audio)
	if err != nil {
		return err
	}
	// Use hex string as filename
	audioFileHashSum := hex.EncodeToString(h.Sum(nil))

	// Save data to DB

	// Save audio to store
	LocStore = fsstore.NewFSStore()
	LocStore.Configure("/tmp/stt/store")
	LocStore.Save(string(audioFileHashSum), audio.Audio)

	// Call remote STT neural network service
	text, err := clients.ReqSTT(audio.Audio)
	if err != nil {
		return err
	}

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