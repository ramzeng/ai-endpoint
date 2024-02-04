package redis

import "github.com/redis/go-redis/v9"
import "context"

var client *redis.Client
var ctx context.Context
