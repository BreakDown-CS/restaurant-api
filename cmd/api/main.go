package main

import (
	"context"
	"log"

	"github.com/BreakDown-CS/restaurant-api/internal/app"
	"github.com/BreakDown-CS/restaurant-api/internal/config"
	"github.com/BreakDown-CS/restaurant-api/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	r, err := app.NewRouter(cfg, pool)
	if err != nil {
		log.Fatal(err)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
