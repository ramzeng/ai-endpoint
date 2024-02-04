package azure

type Config struct {
	Models []string
	Peers  []PeerConfig
}

type PeerConfig struct {
	Key         string
	Endpoint    string
	Deployments []Deployment
	Weight      int64
}
