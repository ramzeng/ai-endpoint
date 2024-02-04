package core

import (
	"fmt"

	"github.com/ramzeng/ai-endpoint/package/config"
	"github.com/ramzeng/ai-endpoint/package/logger"
)

func InitializeLogger() {
	var loggerConfig logger.Config

	err := config.UnmarshalKey("logger", &loggerConfig)

	if err != nil {
		fmt.Print("[Core]: parse logger config failed", err)
		return
	}

	err = logger.Initialize(loggerConfig)

	if err != nil {
		fmt.Print("[Core]: logger initialization failed", err)
	}
}
