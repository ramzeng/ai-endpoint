package balancer

import (
	"sync"
)

func NewBalancer() *Balancer {
	return &Balancer{
		peers: make([]PeerInterface, 0),
	}
}

type Balancer struct {
	peers []PeerInterface
	mutex sync.Mutex
}

func (s *Balancer) AddPeer(peer PeerInterface) {
	s.peers = append(s.peers, peer)
}

func (s *Balancer) GetNextPeer() PeerInterface {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.peers) < 1 {
		return nil
	}

	var totalWeight int64

	selectedPeer := s.peers[0]

	for _, p := range s.peers {
		totalWeight += p.GetEffectiveWeight()
		p.IncreaseCurrentWeight(p.GetEffectiveWeight())

		if p.GetCurrentWeight() > selectedPeer.GetCurrentWeight() {
			selectedPeer = p
		}
	}

	selectedPeer.DecreaseCurrentWeight(totalWeight)

	return selectedPeer
}
