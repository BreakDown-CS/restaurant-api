package app

import (
	"net/http"

	"github.com/BreakDown-CS/restaurant-api/internal/config"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config) (*gin.Engine, error) {

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	v1 := r.Group("/api/v1")
	v1.GET("health", health)

	return r, nil
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    gin.H{"status": "ok"},
	})
}
