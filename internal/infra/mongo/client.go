package mongo

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewClient(ctx context.Context, mongoURI, mongoDBName string) (*Client, error) {
	connectCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(connectCtx, nil)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("MongoDB connected successfully 🚀")

	db := client.Database(mongoDBName)

	return &Client{
		Client:   client,
		Database: db,
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}
