package router

import (
	"log"

	"github.com/gin-gonic/gin"

	broker "github.com/SLANGERES/go-service/internal/Broker"
	handler "github.com/SLANGERES/go-service/internal/Handler"
)

func Router() *gin.Engine {
	// Initialize RabbitMQ connection pool
	if err := broker.InitializePool(); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ pool: %v", err)
	}

	// Initialize log buffer and workers
	handler.InitLogBuffer()

	router := gin.Default()

	router.GET("/health", handler.Healthz)

	router.POST("/log", handler.PostLogging)

	return router
}
