package logger_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

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

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	cfg := logger.Config{
		Level:      logger.DebugLevel,
		Output:     &buf,
		TimeFormat: time.RFC3339,
	}
	log := logger.New(cfg)

	// Test all log levels
	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warning message")
	log.Error("error message")

	output := buf.String()
	expectedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	for _, level := range expectedLevels {
		if !strings.Contains(output, level) {
			t.Errorf("Expected log output to contain %s level", level)
		}
	}
}

func TestLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(logger.Config{Output: &buf})

	// Test with single field
	log.Info("test message", logger.Field{Key: "key", Value: "value"})
	if !strings.Contains(buf.String(), "key=value") {
		t.Error("Expected log output to contain field")
	}

	// Test with multiple fields
	buf.Reset()
	log.Info("test message",
		logger.Field{Key: "key1", Value: "value1"},
		logger.Field{Key: "key2", Value: 123},
	)
	output := buf.String()
	if !strings.Contains(output, "key1=value1") || !strings.Contains(output, "key2=123") {
		t.Error("Expected log output to contain all fields")
	}
}

func TestLoggerWith(t *testing.T) {
	var buf bytes.Buffer
	baseLogger := logger.New(logger.Config{Output: &buf})

	// Create a new logger with additional fields
	childLogger := baseLogger.With(
		logger.Field{Key: "user_id", Value: "123"},
		logger.Field{Key: "role", Value: "admin"},
	)

	// Log with child logger
	childLogger.Info("test message", logger.Field{Key: "action", Value: "login"})

	output := buf.String()
	expectedFields := []string{"user_id=123", "role=admin", "action=login"}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("Expected log output to contain field: %s", field)
		}
	}
}

func TestLoggerWithContext(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(logger.Config{Output: &buf})

	// Create context with various IDs
	ctx := context.Background()
	ctx = logger.WithRequestID(ctx, "req-123")
	ctx = logger.WithUserID(ctx, "user-456")
	ctx = logger.WithSessionID(ctx, "sess-789")

	// Create logger with context
	ctxLogger := log.WithContext(ctx)
	ctxLogger.Info("test message")

	output := buf.String()
	expectedFields := []string{
		"request_id=req-123",
		"user_id=user-456",
		"session_id=sess-789",
	}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("Expected log output to contain field: %s", field)
		}
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	cfg := logger.Config{
		Level:  logger.WarnLevel,
		Output: &buf,
	}
	log := logger.New(cfg)

	// Log messages at different levels
	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warning message")
	log.Error("error message")

	output := buf.String()
	if strings.Contains(output, "DEBUG") || strings.Contains(output, "INFO") {
		t.Error("Expected debug and info messages to be filtered out")
	}
	if !strings.Contains(output, "WARN") || !strings.Contains(output, "ERROR") {
		t.Error("Expected warning and error messages to be included")
	}
}

func TestLoggerFatal(t *testing.T) {
	// Skip this test in normal test runs as it would exit the process
	if os.Getenv("TEST_FATAL") == "" {
		t.Skip("Skipping fatal test")
	}

	var buf bytes.Buffer
	log := logger.New(logger.Config{Output: &buf})
	log.Fatal("fatal message")
}

func TestLoggerConcurrent(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(logger.Config{Output: &buf})

	// Test concurrent logging
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			log.Info("concurrent message", logger.Field{Key: "goroutine", Value: id})
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all messages were logged
	output := buf.String()
	for i := 0; i < 10; i++ {
		if !strings.Contains(output, fmt.Sprintf("goroutine=%d", i)) {
			t.Errorf("Expected log output to contain goroutine %d", i)
		}
	}
}
