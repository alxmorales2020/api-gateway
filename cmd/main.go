package main

import (
	"log"
	"net/http"

	"github.com/alxmorales2020/api-gateway/admin"
	"github.com/alxmorales2020/api-gateway/config"
	"github.com/alxmorales2020/api-gateway/core"
	"github.com/alxmorales2020/api-gateway/plugins/auth"
	"github.com/alxmorales2020/api-gateway/plugins/logging"
	"github.com/alxmorales2020/api-gateway/router"
	"github.com/go-chi/chi/v5"
)

// main initializes the API Gateway, loads the configuration, and starts the HTTP server.
// It sets up the router and listens on port 8080.
// The configuration is loaded from a YAML file named "config.yaml".
// The router is created using the NewRouter function from the router package.
// The server listens for incoming HTTP requests and routes them according to the configuration.
// It also mounts the admin API for managing routes and plugins.
func main() {
	gatewayConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	registerPlugin()

	var store config.RouteStore
	if gatewayConfig.Persistence.MongoDB != nil {
		store, err = config.NewMongoRouteStore(gatewayConfig.Persistence.MongoDB)
		if err != nil {
			log.Fatalf("Error connecting to MongoDB: %v", err)
		}
		log.Println("Loaded route configuration from MongoDB.")
	} else {
		store = config.NewYAMLRouteStore(gatewayConfig.Routes)
		log.Println("Loaded route configuration from config.yaml.")
	}

	// Hot-reloadable app router
	manager, err := router.NewManager(store)
	if err != nil {
		log.Fatalf("router manager: %v", err)
	}

	// Top-level router
	top := chi.NewRouter()

	// Admin API (gets store and a reloader)
	adminHandler := admin.NewAdminHandler(store, manager)
	top.Mount("/admin", adminHandler.Routes())

	top.Mount("/", manager) // app routes served via atomic handler

	top.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("TOP 404: %s %s", r.Method, r.URL.Path)
		http.Error(w, "not found", http.StatusNotFound)
	})

	log.Println("Starting API Gateway on :8080")
	if err := http.ListenAndServe(":8080", top); err != nil {
		log.Fatalf("server: %v", err)
	}
}

// registerPlugin registers a plugin with the core plugin manager.
// It takes a plugin name and a function that returns a new instance of the plugin.
// This function is used to dynamically load plugins at runtime.
// The plugin manager maintains a registry of available plugins and their configurations.
// It allows the API Gateway to extend its functionality by adding new plugins without modifying the core code
func registerPlugin() {
	core.RegisterPlugin("logging", logging.New)
	core.RegisterPlugin("jwt-auth", auth.New)
}
