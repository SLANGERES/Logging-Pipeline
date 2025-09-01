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

// ===== Buffer Settings =====
var (
	logBuffer      chan model.Logging
	batchSize      = 100              // Number of logs per batch
	flushInterval  = 500 * time.Millisecond // Maximum wait before flushing
	numWorkers     = 5                // Number of concurrent workers
	bufferCapacity = 10000        // Buffered channel capacity
)

// Initialize the log buffer and start workers
func InitLogBuffer() {
	logBuffer = make(chan model.Logging, bufferCapacity)
	for i := 0; i < numWorkers; i++ {
		go bufferWorker(i + 1)
	}
}

// Background worker to process buffered logs
func bufferWorker(workerID int) {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	var batch []model.Logging

	for {
		select {
		case logRequest := <-logBuffer:
			batch = append(batch, logRequest)
			if len(batch) >= batchSize {
				sendBatch(batch)
				batch = nil
			}
		case <-ticker.C:
			if len(batch) > 0 {
				sendBatch(batch)
				batch = nil
			}
		}
	}
}

// Send a batch of logs to RabbitMQ efficiently
func sendBatch(batch []model.Logging) {
	// Use the new batch sending function for better performance
	if err := broker.SendBatchToQueue(batch); err != nil {
		slog.Error("Failed to send batch to queue", slog.Any("error", err), slog.Int("batch_size", len(batch)))
		// Fallback to individual sends if batch fails
		for _, logReq := range batch {
			if err := broker.SendDataToQueue(logReq); err != nil {
				slog.Error("Failed to send individual log to queue", slog.Any("error", err), slog.Any("request", logReq))
			}
		}
	} else {
		slog.Debug("Successfully sent batch to queue", slog.Int("batch_size", len(batch)))
	}
}

// ===== Request Validators and Helpers =====
func Validator(ctx *gin.Context, requestData *model.Logging) bool {
	if err := validator.New().Struct(requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return false
	}
	return true
}

func AddSystemGenerateField(ctx *gin.Context, requestData *model.Logging) {
	requestData.Id = uuid.New().String()
	requestData.MetaData.SourceIp = ctx.ClientIP()
	if requestData.MetaData.Region == "" {
		requestData.MetaData.Region = "none"
	}
	requestData.TimeStamp = time.Now().UTC()
}

// ===== HTTP Handlers =====
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

	// Add system fields
	AddSystemGenerateField(ctx, &request)

	// Validate request
	if ok := Validator(ctx, &request); !ok {
		return
	}

	// Log locally
	slog.Info("Incoming log request", slog.Any("request", request))

	// Push to buffer (non-blocking)
	select {
	case logBuffer <- request:
		// Successfully buffered
	default:
		// Buffer full, drop or handle overflow
		slog.Warn("Log buffer full, dropping log", slog.Any("request", request))
	}

	// Respond immediately
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"id":      request.Id,
	})
}
