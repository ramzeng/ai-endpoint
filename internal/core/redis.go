package core

import (
	"github.com/ramzeng/ai-endpoint/package/config"
	"github.com/ramzeng/ai-endpoint/package/logger"
	"github.com/ramzeng/ai-endpoint/package/redis"
	"go.uber.org/zap"
)

func InitializeRedis() {
	var redisConfig redis.Config

	err := config.UnmarshalKey("redis", &redisConfig)

	if err != nil {
		logger.Error("app", "[Core]: parse redis config failed", zap.Error(err))
		return
	}

	err = redis.Initialize(redisConfig)

	if err != nil {
		logger.Error("app", "[Core]: redis initialization failed", zap.Error(err))
	}
}
