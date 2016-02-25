package console

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var (
	gray  = color.New(color.FgHiBlack).SprintFunc()
	white = color.New(color.FgHiWhite).SprintFunc()
)

// Changes the configuration of default console, used for the following functions:
// Trace, Debug, Info, Warning, Error, Panic
func SetDefaultCfg(c Cfg) {
	defaultConsole.cfg = c
}

// A hook intercepts log message and perform certain tasks, like sending email
type Hook interface {
	// Unique Id to identify Hook
	Id() string
	// Action performed by the Hook.
	Action(l Lvl, msg, format string, args ...interface{})
	// Condition that triggers the Hook.
	Match(l Lvl, msg, format string, args ...interface{}) bool
}

// A Writer implements the WriteString method, as the os.File
type Writer interface {
	io.Writer
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
	n := Console{
		mu:    l.mu,
		cfg:   l.cfg,
		w:     l.w,
		hooks: l.hooks,
	}
	n.cfg.prefix = prefix
	return &n
}

// Adds a Hook to the logger.
func (l *Console) Add(h Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks[h.Id()] = h
}

// Release an Hook from the logger.
func (l *Console) Release(h Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()
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
func (l *Console) output(lvl Lvl, format string, args ...interface{}) {
	if l.cfg.Lvl > lvl {
		return
	}
	for i := range args {
		if fn, ok := args[i].(func() string); ok {
			args[i] = fn()
		}
	}
	msg := fmt.Sprintf(format, args...)
	l.executeHooks(lvl, msg, format, args...)
	b := bytes.NewBuffer(nil)
	l.writePrefix(b, lvl)
	if l.cfg.Color {
		msg = white(msg)
	}
	b.WriteString(msg)
	if !strings.HasPrefix(msg, "\n") {
		b.WriteString("\n")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	io.Copy(l.w, b)
}

func (l *Console) executeHooks(lvl Lvl, msg, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, h := range l.hooks {
		if h.Match(lvl, msg, format, args...) {
			h.Action(lvl, msg, format, args...)
		}
	}
}

func (l *Console) writePrefix(b Writer, lvl Lvl) {
	if t := l.cfg.Date.fmt(); t != nil {
		b.WriteString(t(time.Now()) + " ")
	}
	b.WriteString(levels[lvl].GetLabel(l.cfg.Color) + " ")
	if f := l.cfg.File.fmt(); f != nil {
		_, name, line, _ := runtime.Caller(3)
		fl := fmt.Sprintf("[%s:%d]", f(name), line)
		if l.cfg.Color {
			fl = gray(fl)
		}
		b.WriteString(fl + " ")
	}
	if l.cfg.prefix != "" {
		p := l.cfg.prefix
		if l.cfg.Color {
			p = levels[lvl].Color.SprintFunc()(p)
		}
		b.WriteString(p + " ")
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
