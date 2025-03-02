package logger

import (
	"github.com/rs/zerolog"
)

var NopLogger = zerolog.Nop()

type LogContext string

const (
	APISimulation LogContext = "api_simulation"
	SnapshotPoll  LogContext = "snapshot_poll"
)

// CtxLogger wraps zerolog.Logger with context-aware logging capabilities
type CtxLogger struct {
	*zerolog.Logger
	context LogContext
}

// NewCtxLogger creates a new CtxLogger with the specified context
func NewCtxLogger(baseLogger *zerolog.Logger, context LogContext) *CtxLogger {
	return &CtxLogger{
		Logger:  baseLogger,
		context: context,
	}
}

// shouldLog determines if logging should occur based on the context
func (cl *CtxLogger) shouldLog() bool {
	switch cl.context {
	case APISimulation:
		return false // Don't log for API simulation
	case SnapshotPoll:
		return true // Always log for snapshot polling
	default:
		return true // Default to logging
	}
}

// Info logs an info message if the context allows
func (cl *CtxLogger) Info() *zerolog.Event {
	if !cl.shouldLog() {
		return NopLogger.Info()
	}
	return cl.Logger.Info()
}

// Debug logs a debug message if the context allows
func (cl *CtxLogger) Debug() *zerolog.Event {
	if !cl.shouldLog() {
		return NopLogger.Debug()
	}
	return cl.Logger.Debug()
}

// Error logs an error message if the context allows
func (cl *CtxLogger) Error() *zerolog.Event {
	// if !cl.shouldLog() {
	// 	return NopLogger.Error()
	// }
	return cl.Logger.Error()
}

// Warn logs a warning message if the context allows
func (cl *CtxLogger) Warn() *zerolog.Event {
	if !cl.shouldLog() {
		return NopLogger.Warn()
	}
	return cl.Logger.Warn()
}
