package main

import (
	"github.com/AidlyTeam/Aidly-Backend/internal/app"
	"github.com/AidlyTeam/Aidly-Backend/internal/config"
)

// @title API Service
// @description API Service for Aidly
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name session_id
func main() {
	// Setting all the configs and starting the app.
	cfg, err := config.Init("./config")
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}
