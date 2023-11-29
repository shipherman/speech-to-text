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
	ID      int32
	Email   string
	Name    string
	Surname string
}

// Store for audio files
type Store interface {
	Configure() error
	Save() error
	Delete() error
}

// Structure to provide username throuh context to handlers
type UserCtxKey struct{}
