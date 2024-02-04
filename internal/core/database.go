package core

import (
	"github.com/ramzeng/ai-endpoint/package/config"
	"github.com/ramzeng/ai-endpoint/package/database"
	"github.com/ramzeng/ai-endpoint/package/logger"
	"go.uber.org/zap"
)

func InitializeDatabase() {
	var databaseConfig database.Config

	err := config.UnmarshalKey("database", &databaseConfig)

	if err != nil {
		logger.Error("app", "[Core]: parse database config failed", zap.Error(err))
		return
	}

	err = database.Initialize(databaseConfig)

	if err != nil {
		logger.Error("app", "[Core]: database initialization failed", zap.Error(err))
	}
}
