package azure

type Config struct {
	Models   []string
	Backends []BackendConfig
}

type BackendConfig struct {
	Key         string
	Endpoint    string
	Deployments []Deployment
	Weight      int64
}
