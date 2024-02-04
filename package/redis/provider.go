package redis

import "github.com/redis/go-redis/v9"
import "context"

func Initialize(config Config) error {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx = context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}
