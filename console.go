// Console package implements multi priority logger
package console

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
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

type Hook interface {
	// Action performed by the Hook.
	Action(l LogLevel, fileline, msg string)
	// Condition that triggers the Hook.
	Match(l LogLevel, format string, args ...interface{}) bool
}

// A simple hook that copies messages from a certain level.
type SimpleHook struct {
	LogLevel
	*bytes.Buffer
}

// Matches only the log level.
func (s *SimpleHook) Match(l LogLevel, format string, args ...interface{}) bool {
	return l == s.LogLevel
}

// Writes log content in a buffer.
func (s *SimpleHook) Action(l LogLevel, fileline, msg string) {
	s.WriteString(fmt.Sprintf("%s "+_format, time.Now().Format("2006/01/02 15:04:05"), levels[int(l)], fileline, "", msg))
}

var levels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "PANIC"}

type Logger interface {
	Clone(string) Logger
	Add(Hook)
	Release(Hook)
	Output(int, LogLevel, string, ...interface{})
	Trace(string, ...interface{})
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Panic(string, ...interface{})
}

// Creates a Logger that uses standard output.
func Std(level LogLevel) Logger {
	l := logger{log.New(os.Stdout, "", log.LstdFlags), level, "", nil}
	return &l
}

// Creates a custom Logger.
func New(log *log.Logger, level LogLevel) Logger {
	l := logger{log, level, "", nil}
	return &l
}

type logger struct {
	log    *log.Logger
	level  LogLevel
	prefix string
	hooks  []Hook
}

// Creates a copy of the logger with the given prefix.
func (l *logger) Clone(prefix string) Logger {
	newl := *l
	if prefix != "" {
		newl.prefix = prefix + " "
	}
	return &newl
}

// Adds a Hook to the logger.
func (l *logger) Add(h Hook) {
	l.hooks = append(l.hooks, h)
}

// Release an Hook from the logger.
func (l *logger) Release(h Hook) {
	for i := range l.hooks {
		if l.hooks[i] != h {
			continue
		}
		l.hooks = l.hooks[:i+copy(l.hooks[i:], l.hooks[i+1:])]
	}
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
func (l *logger) Output(depth int, lvl LogLevel, format string, args ...interface{}) {
	if l.level > lvl {
		return
	}
	_, file, line, _ := runtime.Caller(1 + depth)
	fileline := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	msg := fmt.Sprintf(format, args...)
	for _, h := range l.hooks {
		if h.Match(lvl, format, args...) {
			h.Action(lvl, fileline, msg)
		}
	}
	l.log.Output(0, fmt.Sprintf(_format, levels[int(lvl)], fileline, l.prefix, msg))
}

var defaultLogger = &logger{log.New(os.Stdout, "", log.LstdFlags), 2, "", nil}

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
