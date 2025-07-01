package auth

import (
	"errors"
	"net/http"

	"github.com/alxmorales2020/api-gateway/core"
)

// AuthPlugin is a plugin that handles authentication for incoming requests.
type AuthPlugin struct{}

// Name returns the name of the plugin.
func (plugin *AuthPlugin) Name() string {
	return "jwt-auth"
}

// Init initializes the plugin with any necessary configuration.
func (plugin *AuthPlugin) Init(config map[string]interface{}) error {
	// Initialize the plugin with any necessary configuration.
	// For this example, we don't have any specific configuration.
	return nil
}

// Execute checks the JWT token in the request header and validates it.
func (plugin *AuthPlugin) Execute(writer http.ResponseWriter, request *http.Request) error {
	// Extract the JWT token from the Authorization header
	token := request.Header.Get("Authorization")
	if token == "" {
		http.Error(writer, "Unauthorized: No token provided", http.StatusUnauthorized)
		return errors.New("no token provided")
	}

	// Here you would implement the logic to validate the JWT token.
	// For this example, we will just check if the token equals "valid-token".
	if token != "valid-token" {
		http.Error(writer, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return errors.New("invalid token")
	}

	// If the token is valid, allow the request to proceed
	return nil
}

// New creates a new instance of the AuthPlugin.
func New() core.Plugin {
	return &AuthPlugin{}
}
