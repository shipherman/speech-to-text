package client

import (
	"context"
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
