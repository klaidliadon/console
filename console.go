// Console package implements multi priority logger
package console

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Logging levels.
const (
	L_TRACE = iota
	L_DEBUG
	L_INFO
	L_WARN
	L_ERROR
	L_PANIC
)

var levels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "PANIC"}

type Logger interface {
	Clone(string) Logger
	Output(int, int, string, ...interface{})
	Trace(string, ...interface{})
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Panic(string, ...interface{})
}

// Creates a Logger that uses standard output.
func Std(level int) Logger {
	l := logger{log.New(os.Stdout, "", log.LstdFlags), level, ""}
	return &l
}

// Creates a custom Logger.
func New(log *log.Logger, level int) Logger {
	l := logger{log, level, ""}
	return &l
}

type logger struct {
	log    *log.Logger
	level  int
	prefix string
}

// Creates a copy of the logger with the given prefix.
func (l *logger) Clone(prefix string) Logger {
	newl := *l
	if prefix != "" {
		newl.prefix = prefix + " "
	}
	return &newl
}

// Writes the log with TRACE level.
func (l *logger) Trace(format string, args ...interface{}) {
	l.Output(1, L_TRACE, format, args...)
}

// Writes the log with DEBUG level.
func (l *logger) Debug(format string, args ...interface{}) {
	l.Output(1, L_DEBUG, format, args...)
}

// Writes the log with INFO level.
func (l *logger) Info(format string, args ...interface{}) {
	l.Output(1, L_INFO, format, args...)
}

// Writes the log with WARN level.
func (l *logger) Warn(format string, args ...interface{}) {
	l.Output(1, L_WARN, format, args...)
}

// Writes the log with ERROR level.
func (l *logger) Error(format string, args ...interface{}) {
	l.Output(1, L_ERROR, format, args...)
}

// Writes the log with PANIC level.
func (l *logger) Panic(format string, args ...interface{}) {
	l.Output(1, L_PANIC, format, args...)
}

// Writes the log with custom level and depth.
func (l *logger) Output(depth, lvl int, format string, args ...interface{}) {
	if l.level > lvl {
		return
	}
	_, file, line, _ := runtime.Caller(1 + depth)
	var fileLine = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	l.log.Output(0, fmt.Sprintf("[%5s] %-20s - %s%s", levels[lvl], fileLine, l.prefix, fmt.Sprintf(format, args...)))
}

var defaultLogger = &logger{log.New(os.Stdout, "", log.LstdFlags), 2, ""}

// Writes the default log with TRACE level.
func Trace(format string, args ...interface{}) {
	defaultLogger.Output(1, L_TRACE, format, args...)
}

// Writes the default log with DEBUG level.
func Debug(format string, args ...interface{}) {
	defaultLogger.Output(1, L_DEBUG, format, args...)
}

// Writes the default log with INFO level.
func Info(format string, args ...interface{}) {
	defaultLogger.Output(1, L_INFO, format, args...)
}

// Writes the default log with WARN level.
func Warn(format string, args ...interface{}) {
	defaultLogger.Output(1, L_WARN, format, args...)
}

// Writes the default log with ERROR level.
func Error(format string, args ...interface{}) {
	defaultLogger.Output(1, L_ERROR, format, args...)
}

// Writes the default log with PANIC level.
func Panic(format string, args ...interface{}) {
	defaultLogger.Output(1, L_PANIC, format, args...)
}
