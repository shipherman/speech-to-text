package cmd

import (
	"context"
	"log"
	"time"
)

func Run() {
	c, err := NewClient()
	if err != nil {
		log.Println(err)
	}

	t := time.Now().Add(time.Minute * cfg.Timeout)
	ctx, cancel := context.WithDeadline(context.Background(), t)
	defer cancel()
	c.SendRequest(ctx)
}
