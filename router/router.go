package router

import (
	"log"
	"net/http"

	"github.com/alxmorales2020/api-gateway/config"
	"github.com/alxmorales2020/api-gateway/core"
	"github.com/alxmorales2020/api-gateway/proxy"
	"github.com/go-chi/chi/v5"
)

// NewRouter initializes a new Chi router with the provided gateway configuration.
func NewRouter(config *config.GatewayConfig) http.Handler {
	router := chi.NewRouter()

	// Register routes based on the configuration
	for _, route := range config.Routes {
		switch route.Method {
		case "GET":
			router.Get(route.Path, generateHandler(route))
		case "POST":
			router.Post(route.Path, generateHandler(route))
		case "PUT":
			router.Put(route.Path, generateHandler(route))
		case "DELETE":
			router.Delete(route.Path, generateHandler(route))
		default:
			log.Printf("Unsupported method %s for path %s", route.Method, route.Path)
		}
	}
	// Add a default route for health checks
	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})

	log.Println("Gateway listening on :8080")
	return router
}

// generateHandler creates an HTTP handler for a given route configuration.
func generateHandler(route config.RouteConfig) http.HandlerFunc {
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

	upstreamProxy, err := proxy.NewReverseProxy(route.Upstream)
	if err != nil {
		log.Printf("Invalid upstream for route %s: %v", route.Path, err)
		return func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "Bad gateway config", http.StatusBadGateway)
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		for _, plugin := range plugins {
			if err := plugin.Execute(writer, request); err != nil {
				return
			}
		}

		upstreamProxy.ServeHTTP(writer, request)
	}
}
