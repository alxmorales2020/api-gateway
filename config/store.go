package config

type RouteStore interface {
	LoadRoutes() ([]RouteConfig, error)
	SaveRoute(route *RouteConfig) error
	DeleteRoute(id string) error
}
