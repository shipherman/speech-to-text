package cmd

import (
	"context"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
)

func Register(context.Context, *sttservice.RegisterRequest) (*sttservice.RegisterResponse, error) {

	return &sttservice.RegisterResponse{}, nil
}

func Login(context.Context, *sttservice.LoginRequest) (*sttservice.LoginResponse, error) {

	return &sttservice.LoginResponse{}, nil
}
