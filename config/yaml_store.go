package config

import (
	"errors"
	"sync"
)

// YAMLRouteStore implements RouteStore using in-memory route definitions loaded from config.yaml.
type YAMLRouteStore struct {
	mu     sync.RWMutex
	routes []RouteConfig
}

// NewYAMLRouteStore creates a new store backed by in-memory routes.
func NewYAMLRouteStore(routes []RouteConfig) *YAMLRouteStore {
	return &YAMLRouteStore{
		routes: routes,
	}
}

// LoadRoutes returns the current set of routes.
func (s *YAMLRouteStore) LoadRoutes() ([]RouteConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.routes, nil
}

// SaveRoute appends a new route to the in-memory list.
func (s *YAMLRouteStore) SaveRoute(route RouteConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Optional: validate or deduplicate route before adding
	s.routes = append(s.routes, route)
	return nil
}

// DeleteRoute removes a route by exact path match.
func (s *YAMLRouteStore) DeleteRoute(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, r := range s.routes {
		if r.Path == path {
			s.routes = append(s.routes[:i], s.routes[i+1:]...)
			return nil
		}
	}
	return errors.New("route not found")
}
