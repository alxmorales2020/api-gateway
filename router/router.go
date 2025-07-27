package router

import (
	"log"
	"net/http"
	"strings"

	"github.com/alxmorales2020/api-gateway/config"
	"github.com/alxmorales2020/api-gateway/core"
	"github.com/alxmorales2020/api-gateway/proxy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter initializes a new Chi router with the provided gateway configuration.
func NewRouter(routes []config.RouteConfig) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)

	// Register routes based on the configuration
	for _, route := range routes {
		log.Printf("Registering route %s %v → %s", route.Path, route.Methods, route.Upstream)

		isPrefix := strings.HasSuffix(route.Path, "*")
		cleanPath := strings.TrimSuffix(route.Path, "*")
		handler := generateHandler(route, cleanPath, isPrefix)
		if handler == nil {
			log.Printf("No methods defined for path %s — skipping", route.Path)
			continue
		}

		if isPrefix {
			subRouter := chi.NewRouter()
			subRouter.Handle("/*", handler)
			router.Mount(cleanPath, subRouter)
			log.Printf("Mounted path prefix route: %s* → %s", cleanPath, route.Upstream)
		} else {
			for _, method := range route.Methods {
				router.Method(method, route.Path, handler)
				log.Printf("Bound %s %s → %s", method, route.Path, route.Upstream)
			}
		}
	}

	// Add a default route for health checks
	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ROUTE NOT FOUND: %s %s", r.Method, r.URL.Path)
		http.Error(w, "Route not found", http.StatusNotFound)
	})

	log.Println("Gateway router initialized.")
	return router
}

// generateHandler creates an HTTP handler for a given route configuration.
func generateHandler(route config.RouteConfig, prefix string, strip bool) http.HandlerFunc {
	plugins := []core.Plugin{}
	for _, name := range route.Plugins {
		plugin := core.GetPlugin(name)
		if plugin == nil {
			log.Printf("Plugin not found: %s", name)
			continue
		}
		_ = plugin.Init(nil)
		plugins = append(plugins, plugin)
	}

	proxyHandler, err := proxy.NewReverseProxy(route.Upstream, prefixIf(strip, prefix))
	if err != nil {
		log.Printf("Proxy error for %s: %v", route.Path, err)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad gateway config", http.StatusBadGateway)
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		recorder := core.NewResponseRecorder(writer)

		for _, plugin := range plugins {
			if err := plugin.Execute(recorder, request); err != nil {
				return
			}
		}

		proxyHandler.ServeHTTP(recorder, request)
	}
}

func prefixIf(condition bool, prefix string) string {
	if condition {
		return prefix
	}
	return ""
}
