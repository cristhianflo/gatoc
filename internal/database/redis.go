// write the logic to initialize redis
package database

import (
	"context"

	"github.com/bachacode/gatoc/internal/config"
	"github.com/go-redis/redis/v8"
)

func NewRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	options, err := redis.ParseURL(cfg.Url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
