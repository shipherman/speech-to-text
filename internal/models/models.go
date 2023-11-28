package models

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
