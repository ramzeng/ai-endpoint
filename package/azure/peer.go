package azure

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"sync"

	error2 "github.com/ramzeng/ai-endpoint/package/error"
	"github.com/ramzeng/ai-endpoint/package/toolkit"
	"go.uber.org/zap"
)

type Deployment struct {
	Name    string
	Model   string
	Version string
}

type Peer struct {
	Key             string
	Endpoint        *url.URL
	Deployments     []Deployment
	ReverseProxy    *httputil.ReverseProxy
	Weight          int64
	CurrentWeight   int64
	EffectiveWeight int64
	logger          *zap.Logger
	mutex           sync.Mutex
}

func (b *Peer) getMaskedKey() string {
	return toolkit.MaskString(b.Key, 0.7)
}

func (b *Peer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	b.ReverseProxy.ServeHTTP(writer, request)
}

func (b *Peer) InitializeReverseProxy() {
	b.ReverseProxy = &httputil.ReverseProxy{
		Transport:  DefaultTransport,
		Director:   b.Director(),
		BufferPool: toolkit.NewBytesBufferPool(32 * 1024),
		ModifyResponse: func(response *http.Response) error {
			if response.StatusCode == http.StatusOK {
				if b.EffectiveWeight < b.Weight {
					b.AddEffectiveWeight(1)
				}
			}

			return nil
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			b.logger.Error(
				"[Azure]: OpenAI proxy request error",
				zap.String("event", "azure_openai_proxy_request_error"),
				zap.String("key", b.getMaskedKey()),
				zap.String("model", request.Header.Get("X-OpenAI-Model")),
				zap.String("version", request.URL.Query().Get("api-version")),
				zap.String("host", request.URL.Host),
				zap.String("path", request.URL.Path),
				zap.String("request_id", request.Header.Get("X-Request-Id")),
				zap.Error(err),
			)

			if errors.Is(err, context.Canceled) {
				writer.WriteHeader(error2.ClientClosedRequest)
			}

			if val, ok := err.(net.Error); ok {
				if val.Timeout() {
					writer.WriteHeader(http.StatusGatewayTimeout)
				}
			}

			b.AddEffectiveWeight(-b.CurrentWeight / 2)
		},
	}
}

func (b *Peer) HasOpenAIModelCapability(model string) bool {
	for _, deployment := range b.Deployments {
		if deployment.Model == model {
			return true
		}
	}
	return false
}

func (b *Peer) GetDeploymentByModel(model string) (Deployment, bool) {
	for _, deployment := range b.Deployments {
		if deployment.Model == model {
			return deployment, true
		}
	}

	return Deployment{}, false
}

func (b *Peer) AddEffectiveWeight(delta int64) {
	b.mutex.Lock()

	b.EffectiveWeight += delta

	if b.EffectiveWeight > b.Weight {
		b.EffectiveWeight = b.Weight
	}

	if b.EffectiveWeight < 1 {
		b.EffectiveWeight = 1
	}

	b.mutex.Unlock()
}

func (b *Peer) Director() func(request *http.Request) {
	return func(request *http.Request) {
		request.Header.Set("api-key", b.Key)
		request.Header.Del("Authorization")

		deployment, _ := b.GetDeploymentByModel(request.Header.Get("X-OpenAI-Model"))

		query := request.URL.Query()
		query.Add("api-version", deployment.Version)

		request.Host = b.Endpoint.Host
		request.URL.Host = b.Endpoint.Host
		request.URL.Scheme = b.Endpoint.Scheme
		request.URL.Path = path.Join(fmt.Sprintf("/openai/deployments/%s", deployment.Name), strings.Replace(request.URL.Path, "/v1/", "/", 1))

		request.URL.RawPath = request.URL.EscapedPath()
		request.URL.RawQuery = query.Encode()

		b.logger.Info(
			"[Azure]: OpenAI proxy request constructed",
			zap.String("event", "azure_openai_proxy_request_constructed"),
			zap.String("key", b.getMaskedKey()),
			zap.String("model", deployment.Model),
			zap.String("deployment", deployment.Name),
			zap.String("version", deployment.Version),
			zap.String("host", request.URL.Host),
			zap.String("path", request.URL.Path),
			zap.String("request_id", request.Header.Get("X-Request-Id")),
		)
	}
}
