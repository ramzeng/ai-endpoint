package azure

import (
	"net/url"

	"github.com/ramzeng/ai-endpoint/package/balancer"
	"go.uber.org/zap"
)

func Initialize(azureProxyConfig Config, logger *zap.Logger) {
	balancers = map[string]*balancer.Balancer{}

	for _, model := range azureProxyConfig.Models {
		balancers[model] = balancer.NewBalancer()

		for _, peerConfig := range azureProxyConfig.Peers {
			for _, deployment := range peerConfig.Deployments {
				b, ok := balancers[deployment.Model]

				if !ok {
					continue
				}

				endpoint, _ := url.Parse(peerConfig.Endpoint)

				peer := &Peer{
					Key:         peerConfig.Key,
					Endpoint:    endpoint,
					Deployments: peerConfig.Deployments,
					logger:      logger,
				}

				peer.Weight = peerConfig.Weight
				peer.EffectiveWeight = peerConfig.Weight

				peer.InitializeReverseProxy()

				b.AddPeer(peer)
			}
		}
	}
}

func SelectPeerByModel(model string) *Peer {
	return balancers[model].GetNextPeer().(*Peer)
}
