package console

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var (
	gray  = color.New(color.FgHiBlack).SprintFunc()
	white = color.New(color.FgHiWhite).SprintFunc()
)

// SetDefaultCfg changes the configuration of default console, used by Trace, Debug, Info, Warning, Error, Panic
func SetDefaultCfg(c Cfg) {
	c.validate()
	defaultConsole.cfg = c
}

// Hook intercepts log message and perform certain tasks, like sending email
type Hook interface {
	// Unique Id to identify Hook
	Id() string
	// Action performed by the Hook.
	Action(l Lvl, msg, format string, args ...interface{})
	// Condition that triggers the Hook.
	Match(l Lvl, msg, format string, args ...interface{}) bool
}

// Writer implements the WriteString method, as the os.File
type Writer interface {
	io.Writer
	WriteString(string) (int, error)
}

// New creates a Console.
func New(c Cfg, w Writer) *Console {
	c.validate()
	return &Console{
		mu:    &sync.Mutex{},
		cfg:   c,
		w:     w,
		hooks: make(map[string]Hook),
	}
}

type Console struct {
	mu    *sync.Mutex
	cfg   Cfg
	w     Writer
	hooks map[string]Hook
}

// Clone creates a copy of the console with the given prefix.
func (c Console) Clone(prefix string) *Console {
	c.cfg.prefix = prefix
	c.cfg.validate()
	return &c
}

// Adds a Hook to the logger.
func (c *Console) Add(h Hook) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hooks[h.Id()] = h
}

// Release an Hook from the logger.
func (c *Console) Release(h Hook) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.hooks, h.Id())
}

// Trace writes the console with LvlTrace.
func (c *Console) Trace(format string, args ...interface{}) {
	c.print(LvlTrace, format, args...)
}

// Debug writes the console with LvlDebug.
func (c *Console) Debug(format string, args ...interface{}) {
	c.print(LvlDebug, format, args...)
}

// Info writes the console with LvlInfo.
func (c *Console) Info(format string, args ...interface{}) {
	c.print(LvlInfo, format, args...)
}

// Warn writes the console with LvlWarn.
func (c *Console) Warn(format string, args ...interface{}) {
	c.print(LvlWarn, format, args...)
}

// Error writes the console with LvlError.
func (c *Console) Error(format string, args ...interface{}) {
	c.print(LvlError, format, args...)
}

// Panic writes the console with LvlPanic.
func (c *Console) Panic(format string, args ...interface{}) {
	c.print(LvlPanic, format, args...)
}

// print writes the log with custom level and depth.
func (c *Console) print(lvl Lvl, format string, args ...interface{}) {
	if c.cfg.Lvl > lvl {
		return
	}
	for i := range args {
		if fn, ok := args[i].(func() string); ok {
			args[i] = fn()
		}
	}
	msg := fmt.Sprintf(format, args...)
	c.executeHooks(lvl, msg, format, args...)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.w.WriteString(fmt.Sprintf(c.cfg.fmt, c.cfg.args(lvl, msg)...))
	if !strings.HasSuffix(msg, "\n") {
		c.w.Write([]byte{'\n'})
	}
}

func (c *Console) executeHooks(lvl Lvl, msg, format string, args ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, h := range c.hooks {
		if h.Match(lvl, msg, format, args...) {
			h.Action(lvl, msg, format, args...)
		}
	}
}

// Defaults is the configuration for standard console
var Defaults = Cfg{
	Color: true,
	Date:  DateHour,
	File:  FileShow,
}

var defaultConsole = Std()

// Std creates a standard Console on `os.Stdout`.
func Std() *Console {
	c := Defaults
	c.validate()
	return &Console{
		mu:    new(sync.Mutex),
		cfg:   c,
		w:     os.Stdout,
		hooks: make(map[string]Hook),
	}
}

// Trace writes the default console with LvlTrace.
func Trace(format string, args ...interface{}) {
	defaultConsole.print(LvlTrace, format, args...)
}

// Debug writes the default console with LvlDebug.
func Debug(format string, args ...interface{}) {
	defaultConsole.print(LvlDebug, format, args...)
}

// Info writes the default console with LvlInfo.
func Info(format string, args ...interface{}) {
	defaultConsole.print(LvlInfo, format, args...)
}

// Warn writes the default console with LvlWarn.
func Warn(format string, args ...interface{}) {
	defaultConsole.print(LvlWarn, format, args...)
}

// Error writes the default console with LvlError.
func Error(format string, args ...interface{}) {
	defaultConsole.print(LvlError, format, args...)
}

// Panic writes the default console with LvlPanic.
func Panic(format string, args ...interface{}) {
	defaultConsole.print(LvlPanic, format, args...)
}
