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

	for _, peerConfig := range azureProxyConfig.Peers {
		for _, deployment := range peerConfig.Deployments {
			serverPool, ok := serverPools[deployment.Model]

			if !ok {
				continue
			}

			endpoint, _ := url.Parse(peerConfig.Endpoint)

			peer := &Peer{
				Key:             peerConfig.Key,
				Endpoint:        endpoint,
				Deployments:     peerConfig.Deployments,
				Weight:          peerConfig.Weight,
				EffectiveWeight: peerConfig.Weight,
				logger:          logger,
			}

			peer.InitializeReverseProxy()

			serverPool.AddPeer(peer)
		}
	}
}
