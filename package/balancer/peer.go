package balancer

type Peer interface {
	GetEffectiveWeight() int64
	GetCurrentWeight() int64
	SetCurrentWeight(currentWeight int64)
	IncreaseCurrentWeight(increase int64)
	DecreaseCurrentWeight(decrease int64)
}
