package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ramzeng/ai-endpoint/internal/azure"
	"github.com/ramzeng/ai-endpoint/internal/core"
	"github.com/ramzeng/ai-endpoint/internal/middleware"
	"github.com/ramzeng/ai-endpoint/package/config"
)

func main() {
	core.ReadConfig(os.Getenv("ENDPOINT_CONFIG_PATH"))
	core.Boot()

	gin.SetMode(config.GetString("gin.mode"))

	server := gin.New()

	server.Use(
		middleware.RequestId(),
		middleware.Logger(),
		middleware.Auth(),
		middleware.Limiter(),
		gin.Recovery(),
	)

	azure.RegisterRoutes(server)

	_ = server.Run(":8080")
}
