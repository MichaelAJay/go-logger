package logger

import (
	"io"
	"os"
)

type LoggerFactory struct {
	defaultConfig Config
}

func NewFactory(defaultConfig Config) *LoggerFactory {
	return &LoggerFactory{
		defaultConfig: defaultConfig,
	}
}

// Global factory instance
var DefaultFactory = NewFactory(DefaultConfig)

func (f *LoggerFactory) Console(level Level) Logger {
	cfg := f.defaultConfig
	cfg.Level = level
	cfg.Output = os.Stdout
	return New(cfg)
}

func (f *LoggerFactory) File(filePath string, level Level) (Logger, error) {
	return CreateFileLogger(filePath, level)
}

func (f *LoggerFactory) Custom(cfg Config) Logger {
	return New(cfg)
}

func (f *LoggerFactory) Combined(filePath string, consoleLevel, fileLevel Level) (Logger, error) {
	consoleLogger := f.Console(consoleLevel)

	fileLogger, err := f.File(filePath, fileLevel)
	if err != nil {
		return nil, err
	}

	return MultiLogger(consoleLogger, fileLogger), nil
}

func (f *LoggerFactory) NewWriter(logger Logger, level Level) io.Writer {
	return &logWriter{
		logger: logger,
		level:  level,
	}
}

type logWriter struct {
	logger Logger
	level  Level
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	msg := string(p)

	switch w.level {
	case DebugLevel:
		w.logger.Debug(msg)
	case InfoLevel:
		w.logger.Info(msg)
	case WarnLevel:
		w.logger.Warn(msg)
	case ErrorLevel:
		w.logger.Error(msg)
	case FatalLevel:
		w.logger.Fatal(msg)
	default:
		w.logger.Info(msg)
	}

	return len(p), nil
}
