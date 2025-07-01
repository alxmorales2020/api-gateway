package logging

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alxmorales2020/api-gateway/core"
)

// LoggingPlugin is a plugin that logs request and response details.
type LoggingPlugin struct{}

func (plugin *LoggingPlugin) Name() string {
	return "LoggingPlugin"
}

func (plugin *LoggingPlugin) Init(config map[string]interface{}) error {
	// Initialize the plugin with any necessary configuration.
	// For this example, we don't have any specific configuration.
	return nil
}

func (plugin *LoggingPlugin) Execute(writer http.ResponseWriter, request *http.Request) error {
	startTime := time.Now()

	// Log the incoming request details
	fmt.Printf("Received %s request for %s\n", request.Method, request.URL.Path)

	// Call the next handler in the chain (if any)
	// Here we just simulate a response for demonstration purposes
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Request processed"))

	// Log the response details
	duration := time.Since(startTime)
	fmt.Printf("Processed %s request for %s in %v\n", request.Method, request.URL.Path, duration)

	return nil
}

// Register the plugin with the core plugin manager
func New() core.Plugin {
	return &LoggingPlugin{}
}
