package main

import (
	"gateway-service/config"
	"gateway-service/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	routes := setupRoutes()
	routes.Run(cfg.AppPort)
}

func setupRoutes() *routes.Routes {
	return &routes.Routes{}
}
