package router

import (
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/alxmorales2020/api-gateway/config"
)

type Reloader interface {
	Reload() error
}

type Manager struct {
	store   config.RouteStore
	current atomic.Value // holds http.Handler
}

func NewManager(store config.RouteStore) (*Manager, error) {
	m := &Manager{store: store}
	if err := m.Reload(); err != nil {
		return nil, err
	}
	return m, nil
}

// ServeHTTP lets Manager be used as an http.Handler.
// It delegates to the current router atomically.
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, _ := m.current.Load().(http.Handler)
	h.ServeHTTP(w, r)
}

func (m *Manager) Reload() error {
	routes, err := m.store.LoadRoutes()
	if err != nil {
		return err
	}
	app := buildAppRouter(routes)
	m.current.Store(app)
	log.Printf("Router reloaded with %d route(s).", len(routes))
	return nil
}

// buildAppRouter is your existing NewRouter but returning a chi.Router
// for the app routes only (no /admin here).
func buildAppRouter(routes []config.RouteConfig) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	for _, route := range routes {
		isPrefix := strings.HasSuffix(route.Path, "*")
		cleanPath := strings.TrimSuffix(route.Path, "*")
		handler := generateHandler(route, cleanPath, isPrefix)
		if handler == nil || len(route.Methods) == 0 {
			continue
		}

		if isPrefix {
			sub := chi.NewRouter()
			sub.Handle("/*", handler) // any method
			r.Mount(cleanPath, sub)
		} else {
			for _, method := range route.Methods {
				r.Method(method, route.Path, handler)
			}
		}
	}

	// health
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("pong")) })

	// 404
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Route not found", http.StatusNotFound)
	})

	return r
}
