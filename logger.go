package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Convert Level to string
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Logger defines the interface for logging operations
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
}

// Config holds logger configuration
type Config struct {
	Level      Level
	Output     io.Writer
	TimeFormat string
	Prefix     string
}

// DefaultConfig provides sensible defaults
var DefaultConfig = Config{
	Level:      InfoLevel,
	Output:     os.Stdout,
	TimeFormat: time.RFC3339,
	Prefix:     "",
}

// standardLogger implements Logger using Go's standard log package
type standardLogger struct {
	logger     *log.Logger
	level      Level
	timeFormat string
	fields     []Field
	mu         sync.Mutex
}

func New(cfg Config) Logger {
	if cfg.Output == nil {
		cfg.Output = DefaultConfig.Output
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = DefaultConfig.TimeFormat
	}

	logger := log.New(cfg.Output, cfg.Prefix, log.LstdFlags)

	return &standardLogger{
		logger:     logger,
		level:      cfg.Level,
		timeFormat: cfg.TimeFormat,
		fields:     []Field{},
	}
}

// formatFields converts fields to a string representation
func (l *standardLogger) formatFields(fields []Field) string {
	if len(fields) == 0 {
		return ""
	}

	result := "{"
	for i, field := range fields {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprintf("%s=%v", field.Key, field.Value)
	}
	result += "}"

	return result
}

func (l *standardLogger) log(level Level, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Combie base fields with method fields
	allFields := append(l.fields, fields...)

	// Format the log entry
	timestamp := time.Now().Format(l.timeFormat)
	formattedFields := l.formatFields(allFields)

	// Log entry format: timestamp [LEVEL] message {fields}
	l.logger.Printf("%s [%s] %s %s", timestamp, level.String(), msg, formattedFields)

	// Exit on fatal errors
	if level == FatalLevel {
		os.Exit(1)
	}
}

func (l *standardLogger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields...)
}

func (l *standardLogger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields...)
}

func (l *standardLogger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields...)
}

func (l *standardLogger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields...)
}

func (l *standardLogger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields...)
}

// With returns a new logger with the given fields added
func (l *standardLogger) With(fields ...Field) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &standardLogger{
		logger:     l.logger,
		level:      l.level,
		timeFormat: l.timeFormat,
		fields:     make([]Field, len(l.fields), len(l.fields)+len(fields)),
	}

	copy(newLogger.fields, l.fields)

	newLogger.fields = append(newLogger.fields, fields...)
	return newLogger
}

// WithContext returns a new logger with context values
func (l *standardLogger) WithContext(ctx context.Context) Logger {
	// Start with the current logger's fields
	newFields := make([]Field, len(l.fields))
	copy(newFields, l.fields)

	// Add request ID if available
	if requestID, ok := GetRequestID(ctx); ok {
		newFields = append(newFields, Field{Key: "request_id", Value: requestID})
	}

	// Add user ID if available
	if userID, ok := GetUserID(ctx); ok {
		newFields = append(newFields, Field{Key: "user_id", Value: userID})
	}

	// Add session ID if available
	if sessionID, ok := GetSessionID(ctx); ok {
		newFields = append(newFields, Field{Key: "session_id", Value: sessionID})
	}

	// Create a new logger with all the fields
	newLogger := &standardLogger{
		logger:     l.logger,
		level:      l.level,
		timeFormat: l.timeFormat,
		fields:     newFields,
	}

	return newLogger
}
