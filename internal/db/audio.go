package db

import (
	"context"

	"github.com/shipherman/speech-to-text/gen/ent"
)

// SaveNewAudio creates new tuple in table Audio
func (c *Connector) SaveNewAudio(a ent.Audio, u *ent.User) (*ent.Audio, error) {
	entAudio, err := c.Client.Audio.Create().
		SetHash(a.Hash).
		SetUser(u).
		SetPath(a.Path).
		Save(context.Background())

	if err != nil {
		return entAudio, err
	}

	return entAudio, err
}
