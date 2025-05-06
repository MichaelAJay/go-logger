# Go Logger

A flexible and feature-rich logging package for Go applications. This package provides structured logging capabilities with support for different log levels, context-based logging, and concurrent-safe operations.

## Features

- Multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL)
- Structured logging with key-value fields
- Context-aware logging (request ID, user ID, session ID)
- Thread-safe operations
- Configurable output and formatting
- Support for custom time formats
- Field inheritance through logger chaining

## Installation

```bash
go get github.com/MichaelAJay/go-logger
```

## Quick Start

```go
package main

import (
    "github.com/MichaelAJay/go-logger"
)

func main() {
    // Create a logger with default configuration
    log := logger.New(logger.DefaultConfig)

    // Basic logging
    log.Info("Application started")
    log.Debug("Debug information")
    log.Warn("Warning message")
    log.Error("Error occurred")

    // Structured logging with fields
    log.Info("User action", 
        logger.Field{Key: "user_id", Value: "123"},
        logger.Field{Key: "action", Value: "login"},
    )
}
```

## Configuration

The logger can be configured using the `Config` struct:

```go
cfg := logger.Config{
    Level:      logger.InfoLevel,    // Set minimum log level
    Output:     os.Stdout,           // Set output destination
    TimeFormat: time.RFC3339,        // Set time format
    Prefix:     "myapp",             // Set log prefix
}

log := logger.New(cfg)
```

## Log Levels

The package supports the following log levels (in ascending order):

- `DebugLevel`: Detailed information for debugging
- `InfoLevel`: General operational information
- `WarnLevel`: Warning messages for potentially harmful situations
- `ErrorLevel`: Error events that might still allow the application to continue
- `FatalLevel`: Critical errors that require the application to exit

## Structured Logging

Add structured fields to your log messages:

```go
log.Info("Processing request",
    logger.Field{Key: "request_id", Value: "req-123"},
    logger.Field{Key: "method", Value: "POST"},
    logger.Field{Key: "path", Value: "/api/users"},
)
```

## Context-Aware Logging

The logger can automatically extract and include context information:

```go
ctx := context.Background()
ctx = logger.WithRequestID(ctx, "req-123")
ctx = logger.WithUserID(ctx, "user-456")
ctx = logger.WithSessionID(ctx, "sess-789")

ctxLogger := log.WithContext(ctx)
ctxLogger.Info("Request processed") // Automatically includes request_id, user_id, and session_id
```

## Logger Chaining

Create child loggers with inherited fields:

```go
// Create a base logger with common fields
baseLogger := log.With(
    logger.Field{Key: "service", Value: "auth"},
    logger.Field{Key: "version", Value: "1.0.0"},
)

// Create a child logger with additional fields
userLogger := baseLogger.With(
    logger.Field{Key: "user_id", Value: "123"},
)

// All logs from userLogger will include service, version, and user_id
userLogger.Info("User action")
```

## Thread Safety

The logger is designed to be thread-safe and can be used concurrently:

```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        log.Info("Concurrent operation", 
            logger.Field{Key: "goroutine", Value: id},
        )
    }(i)
}
wg.Wait()
```

## Best Practices

1. **Log Levels**: Use appropriate log levels to categorize messages
2. **Structured Fields**: Include relevant context in structured fields
3. **Context Usage**: Use context-aware logging for request tracing
4. **Logger Chaining**: Create specialized loggers for different components
5. **Error Handling**: Always include error details in error logs
6. **Performance**: Avoid expensive operations in debug logs

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.