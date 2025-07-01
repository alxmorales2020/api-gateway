package core

import (
	"log"
	"net/http"

	"github.com/alxmorales2020/api-gateway/config"
	"github.com/alxmorales2020/api-gateway/core"
)

var pluginRegistry = map[string]func() Plugin{}

// RegisterPlugin registers a plugin constructor function with the plugin registry.
func RegisterPlugin(name string, constructor func() Plugin) {
	pluginRegistry[name] = constructor
}

// GetPlugin retrieves a plugin constructor by name from the registry.
func GetPlugin(name string) Plugin {
	if constructor, exists := pluginRegistry[name]; exists {
		return constructor()
	}
	return nil
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
		_ = plugin.Init(nil) // future: pass plugin-specific config
		plugins = append(plugins, plugin)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		for _, plugin := range plugins {
			if err := plugin.Execute(writer, request); err != nil {
				return // plugin already wrote to response
			}
		}

		// TODO: proxy to upstream
		writer.Write([]byte("Handled route: " + route.Path))
	}
}
