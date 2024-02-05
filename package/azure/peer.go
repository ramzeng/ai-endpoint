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

func (p *Peer) GetEffectiveWeight() int64 {
	return p.EffectiveWeight
}

func (p *Peer) GetCurrentWeight() int64 {
	return p.CurrentWeight
}

func (p *Peer) SetCurrentWeight(currentWeight int64) {
	p.CurrentWeight = currentWeight
}

func (p *Peer) IncreaseCurrentWeight(increase int64) {
	p.CurrentWeight += increase
}

func (p *Peer) DecreaseCurrentWeight(decrease int64) {
	p.CurrentWeight -= decrease
}

func (p *Peer) getMaskedKey() string {
	return toolkit.MaskString(p.Key, 0.7)
}

func (p *Peer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	p.ReverseProxy.ServeHTTP(writer, request)
}

func (p *Peer) InitializeReverseProxy() {
	p.ReverseProxy = &httputil.ReverseProxy{
		Transport:  DefaultTransport,
		Director:   p.Director(),
		BufferPool: toolkit.NewBytesBufferPool(32 * 1024),
		ModifyResponse: func(response *http.Response) error {
			if response.StatusCode == http.StatusOK {
				if p.EffectiveWeight < p.Weight {
					p.AddEffectiveWeight(1)
				}
			}

			return nil
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			p.logger.Error(
				"[Azure]: OpenAI proxy request error",
				zap.String("event", "azure_openai_proxy_request_error"),
				zap.String("key", p.getMaskedKey()),
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

			p.AddEffectiveWeight(-p.CurrentWeight / 2)
		},
	}
}

func (p *Peer) HasOpenAIModelCapability(model string) bool {
	for _, deployment := range p.Deployments {
		if deployment.Model == model {
			return true
		}
	}
	return false
}

func (p *Peer) GetDeploymentByModel(model string) (Deployment, bool) {
	for _, deployment := range p.Deployments {
		if deployment.Model == model {
			return deployment, true
		}
	}

	return Deployment{}, false
}

func (p *Peer) AddEffectiveWeight(delta int64) {
	p.mutex.Lock()

	p.EffectiveWeight += delta

	if p.EffectiveWeight > p.Weight {
		p.EffectiveWeight = p.Weight
	}

	if p.EffectiveWeight < 1 {
		p.EffectiveWeight = 1
	}

	p.mutex.Unlock()
}

func (p *Peer) Director() func(request *http.Request) {
	return func(request *http.Request) {
		request.Header.Set("api-key", p.Key)
		request.Header.Del("Authorization")

		deployment, _ := p.GetDeploymentByModel(request.Header.Get("X-OpenAI-Model"))

		query := request.URL.Query()
		query.Add("api-version", deployment.Version)

		request.Host = p.Endpoint.Host
		request.URL.Host = p.Endpoint.Host
		request.URL.Scheme = p.Endpoint.Scheme
		request.URL.Path = path.Join(fmt.Sprintf("/openai/deployments/%s", deployment.Name), strings.Replace(request.URL.Path, "/v1/", "/", 1))

		request.URL.RawPath = request.URL.EscapedPath()
		request.URL.RawQuery = query.Encode()

		p.logger.Info(
			"[Azure]: OpenAI proxy request constructed",
			zap.String("event", "azure_openai_proxy_request_constructed"),
			zap.String("key", p.getMaskedKey()),
			zap.String("model", deployment.Model),
			zap.String("deployment", deployment.Name),
			zap.String("version", deployment.Version),
			zap.String("host", request.URL.Host),
			zap.String("path", request.URL.Path),
			zap.String("request_id", request.Header.Get("X-Request-Id")),
		)
	}
}
