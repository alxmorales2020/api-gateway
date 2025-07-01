package main

import (
	"log"
	"net/http"

	"github.com/alxmorales2020/api-gateway/config"
	"github.com/alxmorales2020/api-gateway/router"
)

func main() {
	configFile, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	router := router.NewRouter(configFile)

	log.Println("Starting API Gateway on :8080")
	http.ListenAndServe(":8080", router)
}
