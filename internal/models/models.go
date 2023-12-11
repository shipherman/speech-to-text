package models

// Server application configuration
type Config struct {
	Address    string
	DSN        string
	STTAddress string
}

type Audio struct {
	Path string
}

type User struct {
	ID       int32
	Email    string
	Name     string
	Password string
}

// Store interfase for audio files
type Store interface {
	Save(string, []byte) error
	Get(string) ([]byte, error)
	GetStorePath() string
}

// Structure to provide username throuh context to handlers
type UserCtxKey struct{}
