package config

import (
	"errors"
	"sync"

	"github.com/google/uuid"
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
	out := make([]RouteConfig, len(s.routes))
	copy(out, s.routes)
	return out, nil
}

// SaveRoute appends a new route to the in-memory list.
func (s *YAMLRouteStore) SaveRoute(route *RouteConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if route.ID == "" {
		route.ID = uuid.NewString()
	}
	s.routes = append(s.routes, *route)
	return nil
}

// DeleteRoute removes a route by its ID.
func (s *YAMLRouteStore) DeleteRoute(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, r := range s.routes {
		if r.ID == id {
			s.routes = append(s.routes[:i], s.routes[i+1:]...)
			return nil
		}
	}
	return errors.New("route not found")
}
