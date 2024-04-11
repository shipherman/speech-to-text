package db

import (
	"context"
	"fmt"

	"github.com/shipherman/speech-to-text/gen/ent"
)

// SaveNewAudio creates new tuple in table Audio
func (c *Connector) SaveNewAudio(ctx context.Context, audioHash string, audioText string, path string, u *ent.User) (*ent.Audio, error) {
	err := c.Client.Audio.Create().
		SetHash(audioHash).
		SetUser(u).
		SetPath(path).
		SetText(audioText).Exec(ctx)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Audio.Query().All(ctx)
	fmt.Println(res)
	fmt.Println("results ⬆")
	return nil, err
}
