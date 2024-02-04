package core

import (
	"github.com/ramzeng/ai-endpoint/package/azure"
	"github.com/ramzeng/ai-endpoint/package/config"
	"github.com/ramzeng/ai-endpoint/package/logger"
	"go.uber.org/zap"
)

func InitializeAzureProxy() {
	var azureProxyConfig azure.Config

	err := config.UnmarshalKey("azure.openai", &azureProxyConfig)

	if err != nil {
		logger.Error("app", "[Core]: parse azure config failed", zap.Error(err))
		return
	}

	azure.Initialize(azureProxyConfig, logger.Channel("app"))
}
