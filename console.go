// Console package implements multi priority logger
package console

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type LogLevel int

// Logging levels.
const (
	L_TRACE = LogLevel(iota)
	L_DEBUG
	L_INFO
	L_WARN
	L_ERROR
	L_PANIC
)

const _format = "[%5s] %-20s - %s%s"

// A hook intercepts log message and perform certain tasks, like sending email
type Hook interface {
	// Action performed by the Hook.
	Action(l LogLevel, fileline, msg string)
	// Condition that triggers the Hook.
	Match(l LogLevel, format string, args ...interface{}) bool
}

var levels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "PANIC"}

// Creates a Logger that uses standard output.
func Std(level LogLevel) *Logger {
	l := Logger{log.New(os.Stdout, "", log.LstdFlags), level, "", nil}
	return &l
}

// Creates a custom Logger.
func New(log *log.Logger, level LogLevel) *Logger {
	l := Logger{log, level, "", nil}
	return &l
}

type Logger struct {
	log    *log.Logger
	level  LogLevel
	prefix string
	hooks  []Hook
}

// Creates a copy of the logger with the given prefix.
func (l *Logger) Clone(prefix string, cleanHooks bool) *Logger {
	newl := *l
	if cleanHooks {
		newl.hooks = nil
	}
	if prefix != "" {
		newl.prefix = prefix + " "
	}
	return &newl
}

// Adds a Hook to the logger.
func (l *Logger) Add(h Hook) {
	l.hooks = append(l.hooks, h)
}

// Release an Hook from the logger.
func (l *Logger) Release(h Hook) {
	for i, hook := range l.hooks {
		if hook != h {
			continue
		}
		l.hooks[i], l.hooks = l.hooks[len(l.hooks)-1], l.hooks[:len(l.hooks)-1]
	}
}

// Writes the log with TRACE level.
func (l *Logger) Trace(format string, args ...interface{}) {
	l.output(1, L_TRACE, format, args...)
}

// Writes the log with DEBUG level.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.output(1, L_DEBUG, format, args...)
}

// Writes the log with INFO level.
func (l *Logger) Info(format string, args ...interface{}) {
	l.output(1, L_INFO, format, args...)
}

// Writes the log with WARN level.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.output(1, L_WARN, format, args...)
}

// Writes the log with ERROR level.
func (l *Logger) Error(format string, args ...interface{}) {
	l.output(1, L_ERROR, format, args...)
}

// Writes the log with PANIC level.
func (l *Logger) Panic(format string, args ...interface{}) {
	l.output(1, L_PANIC, format, args...)
}

// Writes the log with custom level and depth.
func (l *Logger) output(depth int, lvl LogLevel, format string, args ...interface{}) {
	if l.level > lvl {
		return
	}
	_, file, line, _ := runtime.Caller(1 + depth)
	fileline := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	for i := range args {
		if fn, ok := args[i].(func() string); ok {
			args[i] = fn()
		}
	}
	msg := fmt.Sprintf(format, args...)
	for _, h := range l.hooks {
		if h.Match(lvl, format, args...) {
			h.Action(lvl, fileline, msg)
		}
	}
	l.log.Output(0, fmt.Sprintf(_format, levels[int(lvl)], fileline, l.prefix, msg))
}

var defaultLogger = &Logger{log.New(os.Stdout, "", log.LstdFlags), 2, "", nil}

// Writes the default log with TRACE level.
func Trace(format string, args ...interface{}) {
	defaultLogger.output(1, L_TRACE, format, args...)
}

// Writes the default log with DEBUG level.
func Debug(format string, args ...interface{}) {
	defaultLogger.output(1, L_DEBUG, format, args...)
}

// Writes the default log with INFO level.
func Info(format string, args ...interface{}) {
	defaultLogger.output(1, L_INFO, format, args...)
}

// Writes the default log with WARN level.
func Warn(format string, args ...interface{}) {
	defaultLogger.output(1, L_WARN, format, args...)
}

// Writes the default log with ERROR level.
func Error(format string, args ...interface{}) {
	defaultLogger.output(1, L_ERROR, format, args...)
}

// Writes the default log with PANIC level.
func Panic(format string, args ...interface{}) {
	defaultLogger.output(1, L_PANIC, format, args...)
}
