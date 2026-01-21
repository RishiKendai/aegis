package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Client struct {
	*redis.Client
}

func NewClient(ctx context.Context, host, password string, db int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Info().Msg("Redis connected successfully 🚀")

	return &Client{
		Client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
