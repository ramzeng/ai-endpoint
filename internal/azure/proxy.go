package azure

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ramzeng/ai-endpoint/package/azure"
	"github.com/ramzeng/ai-endpoint/package/error"
	"github.com/ramzeng/ai-endpoint/package/logger"
	"github.com/ramzeng/ai-endpoint/package/toolkit"
	"go.uber.org/zap"
)

type ProxyRequest struct {
	Model string `binding:"required,oneof=gpt-4-32k gpt-4 gpt-3.5-turbo"`
}

func Proxy(c *gin.Context) {
	body, err := toolkit.ReadRequestBody(c.Request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": error.BadRequest("please check the request body format."),
		})
		return
	}

	var request ProxyRequest

	if err := binding.JSON.BindBody(body, &request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.BadRequest("please check the model parameter."),
		})
		return
	}

	peer := azure.SelectPeerByModel(request.Model)

	if peer == nil {
		logger.Error(
			"app",
			fmt.Sprintf("[Azure]: failed to get %s peer service", request.Model),
			zap.String("event", "select_azure_proxy_peer_error"),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": error.InternalServerError(),
		})
		return
	}

	c.Request.Header.Set("X-OpenAI-Model", request.Model)

	peer.ServeHTTP(c.Writer, c.Request)

	if c.Writer.Header().Get("Content-Type") == "text/event-stream" {
		if _, err := c.Writer.Write([]byte{'\n'}); err != nil {
			logger.Error(
				"app",
				"[Azure]: failed to write SSE newline",
				zap.String("event", "azure_proxy_sse_error"),
				zap.Error(err),
			)
		}
	}
}
