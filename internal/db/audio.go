package db

import (
	"context"

	"github.com/shipherman/speech-to-text/gen/ent"
)

// SaveNewAudio creates new tuple in table Audio
func (c *Connector) SaveNewAudio(ctx context.Context, audioHash string, audioText string, u *ent.User) (*ent.Audio, error) {
	entAudio, err := c.Client.Audio.Create().
		SetHash(audioHash).
		SetUser(u).
		SetText(audioText).
		Save(ctx)
	if err != nil {
		return entAudio, err
	}

	return entAudio, err
}
