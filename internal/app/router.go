// Package app wires the HTTP router: middleware, routes and handlers.
package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BreakDown-CS/restaurant-api/internal/config"
	"github.com/BreakDown-CS/restaurant-api/internal/features/menus"
	"github.com/BreakDown-CS/restaurant-api/internal/validation"
)

func NewRouter(cfg *config.Config, pool *pgxpool.Pool) (*gin.Engine, error) {

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       12 * time.Hour,
	}))

	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	v1 := r.Group("/api/v1")
	v1.GET("health", health(pool))
	validation.Setup()

	menuHandler := menus.NewHandler(menus.NewUsecase(menus.NewRepository(pool)))
	menuHandler.RegisterRoutes(v1)

	return r, nil
}

func health(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "Database unavailable",
				"data":    gin.H{"status": "down", "database": "down"},
			})

			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "OK",
			"data":    gin.H{"status": "ok", "database": "up"},
		})
	}
}
