package logger

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

var logger hclog.Logger

// Use the init process to set the global logger.
// It is expected to be initialized when the plugin starts
// and you need to import the package in the proper order.
func init() {
	level := os.Getenv("YATAS_LOG")
	if level == "" {
		// Do not emit logs by default
		level = "error"
	}

	logger = hclog.New(&hclog.LoggerOptions{
		Level:                    hclog.LevelFromString(level),
		Output:                   os.Stderr,
		TimeFormat:               "15:04:05",
		IncludeLocation:          true,
		AdditionalLocationOffset: 1,
		Color:                    hclog.AutoColor,
		ColorHeaderOnly:          true,
	})
}

// Logger returns hcl.Logger as it is
func Logger() hclog.Logger {
	return logger
}

// Trace emits a message at the TRACE level
func Trace(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Trace(msg, args...)
}

// Debug emits a message at the DEBUG level
func Debug(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Debug(msg, args...)
}

// Info emits a message at the INFO level
func Info(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Info(msg, args...)
}

// Warn emits a message at the WARN level
func Warn(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Warn(msg, args...)
}

// Error emits a message at the ERROR level
func Error(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Error(msg, args...)
}

// SetLogger sets the global logger to the provided logger
func SetLogger(l hclog.Logger) {
	logger = l
}
