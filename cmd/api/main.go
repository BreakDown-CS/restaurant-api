package main

import (
	"log"

	"github.com/BreakDown-CS/restaurant-api/internal/app"
	"github.com/BreakDown-CS/restaurant-api/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r, err := app.NewRouter(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
