package http

import (
	"ratskeller-server/data/config"
	"ratskeller-server/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/nacl/auth"
)

func NewRouter(cfg config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.SecurityHeaders())

	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 8,
		HttpOnly: true,
		Secure:   cfg.IsProduction(),
		SameSite: 2,
	})

	router.Use(sessions.Sessions("rk_session", store))

	contentStore := content.NewFileStore(cfg.ConfigDir)

	public.RegisterRoutes(router, cfg, contentStore)

	V1 := router.Group("/API/V1")
	{
		V1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "ratskeller-server",
			})
		})

		auth.RegisterRotes(V1)

		V1.GET("/content/size", contentStore.GetSite)
		V1.GET("/content/opening-hours", contentStore.GetOpeningHours)

		protected := V1.Group("")
		protected.Use(auth.RequireAuth())
		{
			admin.RegisterRoutes(protected, contentStore)
			analytics.RegisterRoutes(protected)
		}
	}

	return router
}
