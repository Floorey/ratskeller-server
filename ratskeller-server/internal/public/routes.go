package public

import (
	"net/http"
	"path/filepath"
	"ratskeller-server/data/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, cfg config.Config, _ *content.FileStore) {
	router.Static("/asserts", filepath.Join(cfg.PublicDir, "asserts"))

	router.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(cfg.PublicDir, "index"))
	})

	router.GET("/admin", func(c *gin.Context) {
		c.File(filepath.Join(cfg.PublicDir, "index.html"))
	})

	router.Static("/admin/asserts", filepath.Join(cfg.AdminDir, "assert"))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not_found",
		})
	})
}
