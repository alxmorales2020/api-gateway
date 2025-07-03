package main

import (
	"log"
	"net/http"

	"github.com/alxmorales2020/api-gateway/config"
	"github.com/alxmorales2020/api-gateway/core"
	"github.com/alxmorales2020/api-gateway/plugins/auth"
	"github.com/alxmorales2020/api-gateway/plugins/logging"
	"github.com/alxmorales2020/api-gateway/router"
)

// main initializes the API Gateway, loads the configuration, and starts the HTTP server.
// It sets up the router and listens on port 8080.
// The configuration is loaded from a YAML file named "config.yaml".
// The router is created using the NewRouter function from the router package.
// The server listens for incoming HTTP requests and routes them according to the configuration.
func main() {
	configFile, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	registerPlugin()

	router := router.NewRouter(configFile)

	log.Println("Starting API Gateway on :8080")
	http.ListenAndServe(":8080", router)
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
