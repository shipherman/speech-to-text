package jwt

import (
	"testing"
	"time"

	"github.com/shipherman/speech-to-text/gen/ent"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	// Proper token
	_, err := NewToken(ent.User{
		ID:       1,
		Email:    "email",
		Login:    "login",
		Password: "password",
	},
		time.Minute*1,
		"secretstring")

	assert.NoError(t, err)
}
