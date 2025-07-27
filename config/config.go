package config

type RouteConfig struct {
	Path        string   `yaml:"path" bson:"path"`
	Methods     []string `yaml:"methods" bson:"methods"`
	Upstream    string   `yaml:"upstream" bson:"upstream"`
	StripPrefix bool     `yaml:"strip_prefix,omitempty" bson:"strip_prefix,omitempty"`
	Plugins     []string `yaml:"plugins,omitempty" bson:"plugins,omitempty"`
}

type GatewayConfig struct {
	Persistence PersistenceConfig `yaml:"persistence"`
	Routes      []RouteConfig     `yaml:"routes"`
}

type PersistenceConfig struct {
	MongoDB *MongoDBConfig `yaml:"mongodb"`
}

type MongoDBConfig struct {
	URI        string `yaml:"uri"`
	Database   string `yaml:"database"`   // default: apigateway
	Collection string `yaml:"collection"` // default: routes
	Username   string `yaml:"username"`   // optional
	Password   string `yaml:"password"`   // optional
}
