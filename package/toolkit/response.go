package toolkit

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type ResponseWriter struct {
	gin.ResponseWriter
	Buffer bytes.Buffer
}

func (writer *ResponseWriter) Write(b []byte) (int, error) {
	writer.Buffer.Write(b)
	return writer.ResponseWriter.Write(b)
}

func NewResponseWriter(writer gin.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: writer,
		Buffer:         bytes.Buffer{},
	}
}
