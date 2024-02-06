package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type LogLevel string

const (
	// Define log level constants.
	DebugLevel LogLevel = "DEBUG"
	InfoLevel  LogLevel = "INFO"
)

// Logger interface that abstracts logging functions.
type Logger interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

// NewLogger creates a new instance of a Logger with specified log level.
func NewLogger(level LogLevel) Logger {
	var handlerOptions slog.HandlerOptions

	// Enable source file logging.
	handlerOptions.AddSource = true

	// Setup the logger with a JSON handler and configure log level.
	var handler slog.Handler
	switch level {
	case DebugLevel:
		handlerOptions.Level = slog.LevelDebug
	case InfoLevel:
		handlerOptions.Level = slog.LevelInfo
	default:
		handlerOptions.Level = slog.LevelInfo // Default to Info level if unspecified.
	}

	handler = slog.NewTextHandler(os.Stdout, &handlerOptions)
	logger := slog.New(handler)

	return &slogLogger{logger: logger}
}

// slogLogger is an implementation of Logger interface using slog.
type slogLogger struct {
	logger *slog.Logger
}

// Debugf logs a debug message.
func (l *slogLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// Errorf logs an error message.
func (l *slogLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

// Infof logs an informational message.
func (l *slogLogger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Warnf logs a warning message.
func (l *slogLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}
