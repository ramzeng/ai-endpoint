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

	"github.com/ramzeng/ai-endpoint/package/balancer"
	error2 "github.com/ramzeng/ai-endpoint/package/error"
	"github.com/ramzeng/ai-endpoint/package/toolkit"
	"go.uber.org/zap"
)

type Deployment struct {
	Name     string
	Model    string
	Version  string
	IsOpenAI bool
}

type Peer struct {
	balancer.Peer
	Key          string
	Endpoint     *url.URL
	Deployments  []Deployment
	ReverseProxy *httputil.ReverseProxy
	logger       *zap.Logger
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
					p.IncreaseEffectiveWeight(1)
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

			var val net.Error
			if errors.As(err, &val) {
				if val.Timeout() {
					writer.WriteHeader(http.StatusGatewayTimeout)
				}
			}

			p.IncreaseEffectiveWeight(-p.CurrentWeight / 2)
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

func (p *Peer) Director() func(request *http.Request) {
	return func(request *http.Request) {
		request.Host = p.Endpoint.Host
		request.URL.Host = p.Endpoint.Host
		request.URL.Scheme = p.Endpoint.Scheme

		query := request.URL.Query()
		deployment, _ := p.GetDeploymentByModel(request.Header.Get("X-OpenAI-Model"))

		if deployment.IsOpenAI {
			// it's openai original model
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Key))
		} else {
			// it's azure deployed model
			request.Header.Set("api-key", p.Key)
			request.Header.Del("Authorization")
			query.Add("api-version", deployment.Version)
			request.URL.Path = path.Join(fmt.Sprintf("/openai/deployments/%s", deployment.Name), strings.Replace(request.URL.Path, "/v1/", "/", 1))
		}

		request.URL.RawPath = request.URL.EscapedPath()
		request.URL.RawQuery = query.Encode()

		p.logger.Info(
			"[Azure]: OpenAI proxy request constructed",
			zap.String("event", "azure_openai_proxy_request_constructed"),
			zap.String("key", p.getMaskedKey()),
			zap.Bool("is_openai", deployment.IsOpenAI),
			zap.String("model", deployment.Model),
			zap.String("deployment", deployment.Name),
			zap.String("version", deployment.Version),
			zap.String("host", request.URL.Host),
			zap.String("path", request.URL.Path),
			zap.String("request_id", request.Header.Get("X-Request-Id")),
		)
	}
}
