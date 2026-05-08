package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/logout", Logout)
		auth.GET("/me", Me)
	}
}
