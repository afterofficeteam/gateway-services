package main

import (
	"gateway-service/config"
	"gateway-service/routes"
	"gateway-service/util/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	go middleware.CleanupOldLimiters()

	routes := setupRoutes()
	routes.Run(cfg.AppPort)
}

func setupRoutes() *routes.Routes {
	return &routes.Routes{}
}
