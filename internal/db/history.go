package db

import (
	"context"

	"github.com/shipherman/speech-to-text/gen/ent"
	"github.com/shipherman/speech-to-text/gen/ent/audio"
	"github.com/shipherman/speech-to-text/gen/ent/user"
)

// Gethistory returns history of reqs for
// provided user email
func (c *Connector) GetHistory(ctx context.Context, email string) ([]*ent.Audio, error) {
	entAudioArr, err := c.Client.Audio.Query().
		Where(audio.HasUserWith(user.EmailEQ(email))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return entAudioArr, nil
}
