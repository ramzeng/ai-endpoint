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
				balancer, ok := balancers[deployment.Model]

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

				balancer.AddPeer(peer)
			}
		}
	}
}

func SelectPeerByModel(model string) *Peer {
	return balancers[model].GetNextPeer().(*Peer)
}
