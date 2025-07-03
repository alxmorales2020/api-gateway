package config

type RouteConfig struct {
	Method   string   `yaml:"method"`
	Path     string   `yaml:"path"`
	Upstream string   `yaml:"upstream"`
	Plugins  []string `yaml:"plugins"`
}

type GatewayConfig struct {
	Routes []RouteConfig `yaml:"routes"`
}
