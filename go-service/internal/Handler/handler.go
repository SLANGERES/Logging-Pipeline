package handler

import (
	"log/slog"
	"net/http"
	"time"

	broker "github.com/SLANGERES/go-service/internal/Broker"
	model "github.com/SLANGERES/go-service/internal/Model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func Healthz(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Home route",
	})
}
func PostLogging(ctx *gin.Context) {
    var request model.Logging

    // Parse request body
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Add system-generated fields
    request.Id = uuid.New().String()
    request.MetaData.SourceIp = ctx.ClientIP()
    if request.MetaData.Region == "" {
        request.MetaData.Region = "none"
    }
    request.TimeStamp = time.Now().UTC()

    // Validate request
    if err := validator.New().Struct(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
        return
    }

    // Log using slog
    slog.Info("Incoming log request", slog.Any("request", request))

    // Send to RabbitMQ
    if err := broker.SendDataToQueue(request); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Message queue unavailable",
            "details": err.Error(),
        })
        return
    }

    // Send success response
    ctx.JSON(http.StatusOK, gin.H{
        "message": "success",
        "id":      request.Id,
    })
}
