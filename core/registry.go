package core

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
