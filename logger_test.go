package logger_test

import (
	"context"
	"testing"

	"github.com/MichaelAJay/go-logger"
)

func TestLogger(t *testing.T) {
	// Create a simple logger
	log := logger.New(logger.DefaultConfig)

	log.Debug("This is a debug message")
	log.Info("This is an info message")
	log.Warn("This is a warning message", logger.Field{Key: "example", Value: 123})

	// Create logger with fields
	userLogger := log.With(
		logger.Field{Key: "user_id", Value: "abc123"},
		logger.Field{Key: "role", Value: "admin"},
	)

	userLogger.Info("User logged in")

	// Create context with request ID
	ctx := logger.WithRequestID(context.Background(), "req-456")

	// Create logger with context
	requestLogger := log.WithContext(ctx)
	requestLogger.Info("Request received")

	// Use factory to create loggers
	factory := logger.NewFactory(logger.DefaultConfig)

	// Create combined logger
	combinedLogger, err := factory.Combined("logs/application.log", logger.InfoLevel, logger.DebugLevel)
	if err != nil {
		t.Fatalf("Failed to create combined logger: %v", err)
	}

	combinedLogger.Info("This goes to both console and file")
}
