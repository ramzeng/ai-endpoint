package middleware

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ramzeng/ai-endpoint/internal/client"
	"github.com/ramzeng/ai-endpoint/package/logger"
	"github.com/ramzeng/ai-endpoint/package/toolkit"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func formatBody(body []byte) (map[string]interface{}, error) {
	formattedBody := map[string]interface{}{}

	if len(body) == 0 {
		return formattedBody, nil
	}

	err := binding.JSON.BindBody(body, &formattedBody)

	if err != nil {
		return nil, err
	}

	return formattedBody, nil
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		body, _ := toolkit.ReadRequestBody(c.Request)
		formattedRequestBody := gjson.ParseBytes(body).Value()
		writer := toolkit.NewResponseWriter(c.Writer)
		c.Writer = writer

		c.Next()

		latency := time.Since(start)

		var requestClient client.Client
		value, exists := c.Get("client")
		if exists {
			requestClient = value.(client.Client)
		}

		var responseBody []byte
		var reader io.Reader

		contentEncoding := c.Writer.Header().Get("Content-Encoding")
		bytesReader := bytes.NewReader(writer.Buffer.Bytes())

		switch contentEncoding {
		case "gzip":
			reader, _ = gzip.NewReader(bytesReader)
		case "deflate":
			reader = flate.NewReader(bytesReader)
		case "br":
			reader = brotli.NewReader(bytesReader)
		}

		if reader != nil {
			responseBody, _ = io.ReadAll(reader)
		} else {
			responseBody, _ = io.ReadAll(bytesReader)
		}

		formattedResponseBody := gjson.ParseBytes(responseBody).Value()

		logger.Info(
			"request",
			"[Logger]: request log record success",
			zap.Time("@timestamp", start),
			zap.Time("created_at", start),
			zap.String("ip", c.ClientIP()),
			zap.Int64("size", c.Request.ContentLength),
			zap.String("request_id", c.Request.Header.Get("X-Request-Id")),
			zap.Float64("request_time", latency.Seconds()),
			zap.String("upstream_addr", c.Request.Host),
			zap.Int("status", c.Writer.Status()),
			zap.String("request", c.Request.Method+" "+c.Request.URL.Path+" "+c.Request.Proto),
			zap.Any("request_query", c.Request.URL.Query()),
			zap.String("request_body", cast.ToString(body)),
			zap.Any("formatted_request_body", formattedRequestBody),
			zap.Any("request_header", c.Request.Header),
			zap.String("uri", c.Request.RequestURI),
			zap.String("domain", c.Request.Host),
			zap.String("method", c.Request.Method),
			zap.String("referer", c.Request.Header.Get("Referer")),
			zap.String("protocol", c.Request.Proto),
			zap.String("ua", c.Request.UserAgent()),
			zap.Uint64("request_client_id", requestClient.Id),
			zap.String("request_client_name", requestClient.Name),
			zap.String("request_client_ips", c.Request.Header.Get("X-Forwarded-For")),
			zap.Int("response_body_size", c.Writer.Size()),
			zap.String("response_body", string(responseBody)),
			zap.Any("formatted_response_body", formattedResponseBody),
			zap.Any("response_header", c.Writer.Header()),
			zap.Float64("latency", latency.Seconds()),
		)
	}
}
