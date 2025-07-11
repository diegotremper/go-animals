// Package infrastructure provides logger initialization and request-specific log enrichment
// for use throughout the application. It integrates with New Relic for observability support.
package infrastructure

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

// InitLogger initializes a slog.Logger.
// It creates a development or production logger based on the APP_ENV environment variable.
func InitLogger() *slog.Logger {
	var handler slog.Handler
	if os.Getenv("APP_ENV") == "development" {
		handler = slog.NewTextHandler(os.Stdout, nil)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}

	return slog.New(handler)
}

// NewTransactionLogger initializes a slog.Logger with context from the current HTTP transaction.
// It returns a child logger with a "request_id" attribute if present.
func NewTransactionLogger(ctx *gin.Context) *slog.Logger {
	logger := InitLogger()

	requestID := ctx.GetHeader("X-Request-Id")
	if requestID != "" {
		logger = logger.With("request_id", requestID)
	}

	return logger
}
