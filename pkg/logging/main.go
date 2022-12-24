package logging

import (
	"fmt"
	"io"
	"os"
)

// LogLevel represents the level at which a log message should be logged
type LogLevel int

const (

	// DebugLogLevel represents the debug log level
	DebugLogLevel LogLevel = iota

	// InfoLogLevel represents the info log level
	InfoLogLevel

	// WarningLogLevel represents the warning log level
	WarningLogLevel

	// ErrorLogLevel represents the error log level
	ErrorLogLevel
)

// Logger is a basic logger that supports different log levels and colors
type Logger struct {
	AppName string
	out io.Writer
}

// NewLogger creates a new instance of a Logger
func NewLogger(AppName string) *Logger {
	
	return &Logger{AppName: AppName}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.log(l.AppName, DebugLogLevel, "\x1b[36m", msg)
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.log(l.AppName,InfoLogLevel, "\x1b[32m", msg)
}

// Warning logs a warning message
func (l *Logger) Warning(msg string) {
	l.log(l.AppName,WarningLogLevel, "\x1b[33m", msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.log(l.AppName,ErrorLogLevel, "\x1b[31m", msg)
}

func (l *Logger) log(pkg string, level LogLevel, color, msg string) {
	if l.out == nil {
		l.out = os.Stdout
	}

	levelStr := ""
	switch level {
	case DebugLogLevel:
		levelStr = "DEBUG"
	case InfoLogLevel:
		levelStr = "INFO"
	case WarningLogLevel:
		levelStr = "WARNING"
	case ErrorLogLevel:
		levelStr = "ERROR"
	}

	fmt.Fprintf(l.out, "%s[%s][%s]%s %s\n", color, levelStr, pkg, "\x1b[0m", msg)
}