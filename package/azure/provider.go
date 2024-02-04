package azure

import (
	"net/url"

	"go.uber.org/zap"
)

func Initialize(azureProxyConfig Config, logger *zap.Logger) {
	serverPools = map[string]*ServerPool{}

	for _, model := range azureProxyConfig.Models {
		serverPools[model] = &ServerPool{
			peers: []*Peer{},
		}
	}

	for _, backendConfig := range azureProxyConfig.Backends {
		for _, deployment := range backendConfig.Deployments {
			serverPool, ok := serverPools[deployment.Model]

			if !ok {
				continue
			}

			endpoint, _ := url.Parse(backendConfig.Endpoint)

			backend := &Peer{
				Key:             backendConfig.Key,
				Endpoint:        endpoint,
				Deployments:     backendConfig.Deployments,
				Weight:          backendConfig.Weight,
				EffectiveWeight: backendConfig.Weight,
				logger:          logger,
			}

			backend.InitializeReverseProxy()

			serverPool.AddPeer(backend)
		}
	}
}
