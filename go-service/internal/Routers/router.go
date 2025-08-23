package router

import (
	"github.com/gin-gonic/gin"

	"github.com/SLANGERES/go-service/internal/Handler"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.GET("/health",handler.Healthz)

	router.POST("/log", handler.PostLogging)

	return router
}
