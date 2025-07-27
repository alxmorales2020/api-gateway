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
	recorder, ok := writer.(*core.ResponseRecorder)
	if !ok {
		fmt.Println("WARNING: ResponseWriter not wrapped")
		return nil
	}

	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime)
		fmt.Printf("[%s] %s %d %dB %v\n",
			request.Method,
			request.URL.Path,
			recorder.StatusCode,
			recorder.Bytes,
			duration,
		)
	}()

	return nil
}

// Register the plugin with the core plugin manager
func New() core.Plugin {
	return &LoggingPlugin{}
}
