package balancer

import (
	"sync"
)

func NewBalancer() *Balancer {
	return &Balancer{
		peers: make([]Peer, 0),
	}
}

type Balancer struct {
	peers []Peer
	mutex sync.Mutex
}

func (s *Balancer) AddPeer(peer Peer) {
	s.peers = append(s.peers, peer)
}

func (s *Balancer) GetNextPeer() Peer {
	if len(s.peers) < 1 {
		return nil
	}

	s.mutex.Lock()

	defer s.mutex.Unlock()

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
