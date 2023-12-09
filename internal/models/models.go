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

// Store for audio files
// Move to separate package?
type Store interface {
	Configure(string) error
	Save(string, []byte) error
	Delete(string) error
}

// Structure to provide username throuh context to handlers
type UserCtxKey struct{}
