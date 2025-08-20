package config

type RouteConfig struct {
	ID          string   `json:"id,omitempty"			bson:"_id,omitempty"			yaml:"-"` // controlled string id
	Path        string   `json:"path"        			bson:"path"        				yaml:"path"`
	Methods     []string `json:"methods"     			bson:"methods"     				yaml:"methods"`
	Upstream    string   `json:"upstream"    			bson:"upstream"    				yaml:"upstream"`
	StripPrefix bool     `json:"strip_prefix,omitempty" bson:"strip_prefix,omitempty" 	yaml:"strip_prefix,omitempty"`
	Plugins     []string `json:"plugins,omitempty"   	bson:"plugins,omitempty"      	yaml:"plugins,omitempty"`
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
