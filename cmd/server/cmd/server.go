package cmd

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
	"github.com/shipherman/speech-to-text/internal/db"
	clients "github.com/shipherman/speech-to-text/internal/lib"
	"github.com/shipherman/speech-to-text/internal/models"
	"github.com/shipherman/speech-to-text/internal/services/auth"
	"github.com/shipherman/speech-to-text/pkg/audioconverter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TranscribeServer struct {
	sttservice.UnimplementedSttServiceServer
	DBClient db.Connector
	Store    models.Store
	auth     auth.Auth
}

// TranscribeAudio provides main functionality of service
// It adds new audio to the pipeline, sends processing statuses to client
// and provides result of transcription to client
func (t *TranscribeServer) TranscribeAudio(
	audio *sttservice.Audio,
	stream sttservice.SttService_TranscribeAudioServer,
) error {
	// Additional
	//	convert *.ogg to *.wav [for telegram bot]

	var response *sttservice.Status

	// Check wave header and if its valid send status ACCEPTED
	_, err := audioconverter.CheckWAVHeader(audio.Audio)
	if err != nil {
		return err
	}

	// Send status ACCEPTED to client
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

	// Execute user email frome jwt
	// get User object from db

	// Call remote STT neural network service
	audioText, err := clients.ReqSTT(audio.Audio)
	if err != nil {
		return err
	}

	// Save data to DB
	// - execute email from auth token -
	email, err := t.auth.GetEmail(stream.Context())
	if err != nil {
		return err
	}
	// - get user obj from db -
	// Get context from the server stream
	ctx := stream.Context()
	user, err := t.DBClient.GetUser(ctx, email)
	if err != nil {
		return err
	}
	// - save to db -
	path := t.Store.GetStorePath() + "/" + audioFileHashSum
	_, err = t.DBClient.SaveNewAudio(ctx, audioFileHashSum, audioText, path, user)
	if err != nil {
		return err
	}

	// Save audio to store
	err = t.Store.Save(audioFileHashSum, audio.Audio)
	if err != nil {
		return err
	}

	// Send status ORDERED to client
	response = &sttservice.Status{Status: sttservice.EnumStatus_STATUS_ORDERED}
	stream.Send(response)

	// Send status DONE to client
	// with text inside
	response = &sttservice.Status{
		Status: sttservice.EnumStatus_STATUS_DONE,
		Text: &sttservice.Text{
			Text: audioText,
			Len:  int32(len(audioText)),
		},
	}
	stream.Send(response)

	return nil
}

// GetHistory returns array of Text fields to a client
func (t *TranscribeServer) GetHistory(
	in *sttservice.User,
	server sttservice.SttService_GetHistoryServer,
) error {
	var text sttservice.Text

	//Check user email
	email, err := t.auth.GetEmail(server.Context())
	if err != nil {
		return err
	}
	if email != in.Email {
		return fmt.Errorf("wrong email provided")
	}

	history, err := t.DBClient.GetHistory(server.Context(), in.Email)
	if err != nil {
		return err
	}

	for _, k := range history {
		text.Text = k.Text
		text.Len = text.GetLen()
		text.Timestamp = k.Timestamp.String()

		err = server.Send(&text)
		if err != nil {
			return err
		}
	}
	return nil
}

// Register registers new user in stt service
func (t *TranscribeServer) Register(ctx context.Context,
	in *sttservice.RegisterRequest,
) (*sttservice.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := t.auth.RegisterNewUser(ctx, in.Username, in.Email, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sttservice.RegisterResponse{UserId: userID}, nil
}

// Login authenticates users
func (t *TranscribeServer) Login(ctx context.Context, in *sttservice.LoginRequest) (*sttservice.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := t.auth.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &sttservice.LoginResponse{Token: token}, nil
}
