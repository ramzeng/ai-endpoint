package azure

import (
	"sync"
)

type ServerPool struct {
	peers []*Peer
	mutex sync.Mutex
}

func (s *ServerPool) AddPeer(peer *Peer) {
	s.peers = append(s.peers, peer)
}

func (s *ServerPool) GetNextPeer() *Peer {
	if len(s.peers) < 1 {
		return nil
	}

	s.mutex.Lock()

	defer s.mutex.Unlock()

	var totalWeight int64

	selectedPeer := s.peers[0]

	for _, p := range s.peers {
		weight := p.EffectiveWeight

		totalWeight += weight
		p.CurrentWeight += weight

		if p.CurrentWeight > selectedPeer.CurrentWeight {
			selectedPeer = p
		}
	}

	selectedPeer.CurrentWeight -= totalWeight

	return selectedPeer
}

func SelectPeerByModel(model string) *Peer {
	return serverPools[model].GetNextPeer()
}
