package db

import (
	"context"

	"github.com/shipherman/speech-to-text/gen/ent"
	"github.com/shipherman/speech-to-text/internal/models"
)

// SaveNewAudio creates new tuple in table Audio
func (c *Connector) SaveNewAudio(audioHash string, store models.Store, u *ent.User) (*ent.Audio, error) {
	entAudio, err := c.Client.Audio.Create().
		SetHash(audioHash).
		SetUser(u).
		SetPath(store.GetStorePath()).
		Save(context.Background())
	if err != nil {
		return entAudio, err
	}

	return entAudio, err
}
