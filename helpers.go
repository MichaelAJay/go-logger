package logger

import (
	"context"
	"os"
	"path/filepath"
)

// Global logger instance
var defaultLogger Logger

func init() {
	defaultLogger = New(DefaultConfig)
}

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

func GetDefaultLogger() Logger {
	return defaultLogger
}

func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	defaultLogger.Fatal(msg, fields...)
}

func WithContext(ctx context.Context) Logger {
	return defaultLogger.WithContext(ctx)
}

func CreateFileLogger(filePath string, level Level) (Logger, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// Create logger with file output
	cfg := Config{
		Level:      level,
		Output:     file,
		TimeFormat: DefaultConfig.TimeFormat,
	}

	return New(cfg), nil
}

func MultiLogger(loggers ...Logger) Logger {
	return &multiLogger{loggers: loggers}
}

type multiLogger struct {
	loggers []Logger
}

func (m *multiLogger) Debug(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Debug(msg, fields...)
	}
}

func (m *multiLogger) Info(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Info(msg, fields...)
	}
}

func (m *multiLogger) Warn(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Warn(msg, fields...)
	}
}

func (m *multiLogger) Error(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Error(msg, fields...)
	}
}

func (m *multiLogger) Fatal(msg string, fields ...Field) {
	// Only the last logger will exit
	for i, logger := range m.loggers {
		if i == len(m.loggers)-1 {
			logger.Fatal(msg, fields...)
		} else {
			logger.Error(msg, fields...)
		}
	}
}

func (m *multiLogger) With(fields ...Field) Logger {
	newLoggers := make([]Logger, len(m.loggers))
	for i, logger := range m.loggers {
		newLoggers[i] = logger.With(fields...)
	}
	return &multiLogger{loggers: newLoggers}
}

func (m *multiLogger) WithContext(ctx context.Context) Logger {
	newLoggers := make([]Logger, len(m.loggers))
	for i, logger := range m.loggers {
		newLoggers[i] = logger.WithContext(ctx)
	}
	return &multiLogger{loggers: newLoggers}
}
