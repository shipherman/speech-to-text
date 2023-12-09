package db

import (
	"context"
	"fmt"

	"github.com/shipherman/speech-to-text/gen/ent"
	"github.com/shipherman/speech-to-text/gen/ent/user"
)

type Connector struct {
	*ent.Client
}

// CreateUser adds new user to DB
func (c *Connector) CreateUser(ctx context.Context,
	username string,
	email string,
	password string,
) error {
	_, err := c.User.Create().
		SetLogin(username).
		SetEmail(email).
		SetPassword(password).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil

}

func (c *Connector) Login(ctx context.Context,
	email string,
	password string,
) error {
	entUser, err := c.User.Query().
		Where(user.EmailEQ(email)).
		First(ctx)
	if err != nil {
		return err
	}
	if entUser.Password != password {
		return fmt.Errorf("password is incorrect")
	}

	return nil
}
