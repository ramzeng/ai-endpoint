package balancer

import "sync"

type PeerInterface interface {
	GetEffectiveWeight() int64
	GetCurrentWeight() int64
	SetCurrentWeight(currentWeight int64)
	IncreaseCurrentWeight(increase int64)
	DecreaseCurrentWeight(decrease int64)
	IncreaseEffectiveWeight(increase int64)
}

type Peer struct {
	Weight          int64
	CurrentWeight   int64
	EffectiveWeight int64
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

func (p *Peer) IncreaseEffectiveWeight(delta int64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.EffectiveWeight += delta

	if p.EffectiveWeight > p.Weight {
		p.EffectiveWeight = p.Weight
	}

	if p.EffectiveWeight < 1 {
		p.EffectiveWeight = 1
	}
}
