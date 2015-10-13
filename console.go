package console

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// A hook intercepts log message and perform certain tasks, like sending email
type Hook interface {
	// Unique Id to identify Hook
	Id() string
	// Action performed by the Hook.
	Action(l LogLvl, msg string)
	// Condition that triggers the Hook.
	Match(l LogLvl, format string, args ...interface{}) bool
}

// A Writer implements the WriteString method, as the os.File
type Writer interface {
	WriteString(string) (int, error)
}

// Creates a Console.
func New(cfg Cfg, w Writer) *Console {
	return &Console{&sync.Mutex{}, cfg, w, make(map[string]Hook)}
}

type Console struct {
	mu    *sync.Mutex
	cfg   Cfg
	w     Writer
	hooks map[string]Hook
}

// Creates a copy of the logger with the given prefix.
func (l *Console) Clone(prefix string) *Console {
	n := Console{&sync.Mutex{}, l.cfg, l.w, l.hooks}
	n.cfg.prefix = prefix
	return &n
}

// Adds a Hook to the logger.
func (l *Console) Add(h Hook) {
	l.hooks[h.Id()] = h
}

// Release an Hook from the logger.
func (l *Console) Release(h Hook) {
	delete(l.hooks, h.Id())
}

// Writes the log with TRACE level.
func (l *Console) Trace(format string, args ...interface{}) {
	l.output(LvlTrace, format, args...)
}

// Writes the log with DEBUG level.
func (l *Console) Debug(format string, args ...interface{}) {
	l.output(LvlDebug, format, args...)
}

// Writes the log with INFO level.
func (l *Console) Info(format string, args ...interface{}) {
	l.output(LvlInfo, format, args...)
}

// Writes the log with WARN level.
func (l *Console) Warn(format string, args ...interface{}) {
	l.output(LvlWarn, format, args...)
}

// Writes the log with ERROR level.
func (l *Console) Error(format string, args ...interface{}) {
	l.output(LvlError, format, args...)
}

// Writes the log with PANIC level.
func (l *Console) Panic(format string, args ...interface{}) {
	l.output(LvlPanic, format, args...)
}

// Writes the log with custom level and depth.
func (l *Console) output(lvl LogLvl, format string, args ...interface{}) {
	if l.cfg.Lvl > lvl {
		return
	}
	for i := range args {
		if fn, ok := args[i].(func() string); ok {
			args[i] = fn()
		}
	}
	msg := fmt.Sprintf(format, args...)
	for _, h := range l.hooks {
		if h.Match(lvl, format, args...) {
			h.Action(lvl, msg)
		}
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writePrefix(lvl)
	l.w.WriteString(msg)
	if !strings.HasPrefix(msg, "\n") {
		l.w.WriteString("\n")
	}
}

func (l *Console) addSpace() {
	l.w.WriteString(" ")
}

func (l *Console) writePrefix(lvl LogLvl) {
	if t := l.cfg.Date.fmt(); t != nil {
		l.w.WriteString(t(time.Now()))
		l.addSpace()
	}
	l.w.WriteString(listOfLvls[lvl].GetLabel(l.cfg.Color))
	l.addSpace()
	if f := l.cfg.File.fmt(); f != nil {
		_, name, line, _ := runtime.Caller(3)
		l.w.WriteString(fmt.Sprintf("[%s:%d]", f(name), line))
		l.addSpace()
	}
	if l.cfg.prefix != "" {
		l.w.WriteString(l.cfg.prefix)
		l.addSpace()
	}
}

var baseCfg = Cfg{Color: true, Date: DateHour, File: FileShow}
var defaultConsole = Std()

// Creates a standard Console on `os.Stdout`.
func Std() *Console {
	l := Console{&sync.Mutex{}, baseCfg, os.Stdout, make(map[string]Hook)}
	return &l
}

// Writes the default log with a Trace level.
func Trace(format string, args ...interface{}) {
	defaultConsole.output(LvlTrace, format, args...)
}

// Writes the default log with a Debug level.
func Debug(format string, args ...interface{}) {
	defaultConsole.output(LvlDebug, format, args...)
}

// Writes the default log with a Info level.
func Info(format string, args ...interface{}) {
	defaultConsole.output(LvlInfo, format, args...)
}

// Writes the default log with a Warn level.
func Warn(format string, args ...interface{}) {
	defaultConsole.output(LvlWarn, format, args...)
}

// Writes the default log with a Error level.
func Error(format string, args ...interface{}) {
	defaultConsole.output(LvlError, format, args...)
}

// Writes the default log with a Panic level.
func Panic(format string, args ...interface{}) {
	defaultConsole.output(LvlPanic, format, args...)
}
