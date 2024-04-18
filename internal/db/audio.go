package db

import (
	"context"
	"time"

	"github.com/shipherman/speech-to-text/gen/ent"
)

// SaveNewAudio creates new tuple in table Audio
func (c *Connector) SaveNewAudio(ctx context.Context, audioHash string, audioText string, path string, u *ent.User) (*ent.Audio, error) {
	err := c.Client.Audio.Create().
		SetHash(audioHash).
		SetUser(u).
		SetPath(path).
		SetTimestamp(time.Now().UTC()).
		SetText(audioText).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return nil, err
}
