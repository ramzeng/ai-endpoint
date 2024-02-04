package toolkit

import (
	"bytes"
	"io"
	"net/http"
)

func ReadRequestBody(request *http.Request) ([]byte, error) {
	if request.ContentLength == 0 {
		return []byte{}, nil
	}

	buffer := &bytes.Buffer{}

	teaReader := io.TeeReader(request.Body, buffer)

	request.Body = io.NopCloser(buffer)

	return io.ReadAll(teaReader)
}
