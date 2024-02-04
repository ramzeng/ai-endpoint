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

	selectedBackend := s.peers[0]

	for _, b := range s.peers {
		weight := b.EffectiveWeight

		totalWeight += weight
		b.CurrentWeight += weight

		if b.CurrentWeight > selectedBackend.CurrentWeight {
			selectedBackend = b
		}
	}

	selectedBackend.CurrentWeight -= totalWeight

	return selectedBackend
}

func SelectPeerByModel(model string) *Peer {
	return serverPools[model].GetNextPeer()
}
