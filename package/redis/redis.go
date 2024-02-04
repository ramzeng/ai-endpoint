package redis

import (
	"github.com/redis/go-redis/v9"
	"time"
)

func Client() *redis.Client {
	return client
}

func IncrNX(key string, expiration time.Duration) (uint64, error) {
	script := redis.NewScript(`
		local hits = tonumber(redis.call('INCR', KEYS[1]))

		if redis.call('TTL', KEYS[1]) == -1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end

		return hits
	`)

	return script.Run(ctx, client, []string{key}, expiration.Seconds()).Uint64()
}

func Incr(key string) (uint64, error) {
	return client.Incr(ctx, key).Uint64()
}

func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return client.SetNX(ctx, key, value, expiration).Result()
}

func TTL(key string) (time.Duration, error) {
	return client.TTL(ctx, key).Result()
}
