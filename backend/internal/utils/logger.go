package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

// Logger provides structured logging capabilities
type Logger struct {
	level LogLevel
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// NewLogger creates a new logger with the specified level
func NewLogger(level string) *Logger {
	logLevel := INFO
	switch level {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warn":
		logLevel = WARN
	case "error":
		logLevel = ERROR
	}

	return &Logger{
		level: logLevel,
	}
}

// shouldLog checks if a message at the given level should be logged
func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
	}
	return levels[level] >= levels[l.level]
}

// log writes a structured log entry
func (l *Logger) log(level LogLevel, message string, context map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     string(level),
		Message:   message,
		Context:   context,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonBytes))
}

// Debug logs a debug message
func (l *Logger) Debug(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(DEBUG, message, ctx)
}

// Info logs an info message
func (l *Logger) Info(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(INFO, message, ctx)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(WARN, message, ctx)
}

// Error logs an error message
func (l *Logger) Error(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(ERROR, message, ctx)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// WithContext returns a new context map with the given key-value pairs
func WithContext(keyValues ...interface{}) map[string]interface{} {
	ctx := make(map[string]interface{})
	for i := 0; i < len(keyValues)-1; i += 2 {
		key, ok := keyValues[i].(string)
		if ok {
			ctx[key] = keyValues[i+1]
		}
	}
	return ctx
}
