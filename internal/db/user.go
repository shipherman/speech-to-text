package db

import (
	"context"

	"github.com/shipherman/speech-to-text/gen/ent"
	"github.com/shipherman/speech-to-text/gen/ent/user"
)

type Connector struct {
	*ent.Client
}

// CreateUser adds new user to DB
func (c *Connector) SaveUser(ctx context.Context,
	username string,
	email string,
	password string,
) (int64, error) {
	entUser, err := c.User.Create().
		SetLogin(username).
		SetEmail(email).
		SetPassword(password).
		Save(ctx)
	if err != nil {
		return 0, err
	}

	return int64(entUser.ID), nil
}

func (c *Connector) GetUser(ctx context.Context,
	email string,
) (*ent.User, error) {
	entUser, err := c.User.Query().
		Where(user.EmailEQ(email)).
		First(ctx)
	if err != nil {
		return entUser, err
	}

	return entUser, nil
}
