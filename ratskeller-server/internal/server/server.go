package server

import (
	"ratskeller-server/data/config"
	httpserver "ratskeller-server/internal/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    config.Config
	router *gin.Engine
}

func New(cfg config.Config) *Server {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := httpserver.NewRouter(cfg)
}
