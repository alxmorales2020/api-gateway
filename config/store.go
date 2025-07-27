package config

type RouteStore interface {
	LoadRoutes() ([]RouteConfig, error)
	SaveRoute(route RouteConfig) error
	DeleteRoute(path string) error // simple example; can evolve to ID-based
}
