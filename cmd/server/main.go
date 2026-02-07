package main

import (
	"log"

	"github.com/gulmix/Social-Network/internal/config"
	"github.com/gulmix/Social-Network/internal/database"
	"github.com/gulmix/Social-Network/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.InitPostgres(cfg); err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	defer database.ClosePostgres()

	if err := database.InitRedis(cfg); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer database.CloseRedis()

	router := server.SetupRouter(cfg)

	log.Printf("Server starting on http://%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("GraphQL playground available at http://%s:%s/", cfg.Server.Host, cfg.Server.Port)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
