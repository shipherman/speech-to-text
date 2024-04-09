package db

import (
	"context"
	"fmt"

	"github.com/shipherman/speech-to-text/gen/ent"
)

func ConfigureSchema(ctx context.Context, connstring string) error {
	client, err := ent.Open("postgres", connstring)
	if err != nil {
		return err
	}

	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed creating schema resources: %v", err)
	}

	return nil
}
