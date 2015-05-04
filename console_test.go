package console

import (
	"bytes"
	"log"
	"os"
	"testing"
)

// A simple hook that copies in a buffer messages from a certain level.
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
	s.WriteString(msg)
}

// Testing default functions
func TestDefault(t *testing.T) {
	Trace("_")
	Debug("_")
	Info("_")
	Error("_")
	Warn("_")
	Panic("_")
}

// Testing a hook with level error and a function argument
func TestHook(t *testing.T) {
	s := SimpleHook{L_ERROR, bytes.NewBuffer(nil)}
	l := Std(L_TRACE)
	l.Add(&s)
	// message added to hook
	l.Error("message - %s", "args")
	// message to a cloned log ignored
	l.Clone("[prefix]", true).Error("some text %q", "a")
	// ignored for diffent level
	l.Panic("not added")
	l.Warn("not added")
	l.Info("not added")
	l.Debug("not added")
	l.Trace("not added")
	if s.String() != "message - args" {
		t.Error("Unexpected string")
	}
}

func TestFunction(t *testing.T) {
	s := SimpleHook{L_INFO, bytes.NewBuffer(nil)}
	l := Std(L_TRACE)
	l.Add(&s)
	l.Info("%s", func() string { return "x_x_x" })
	if s.String() != "x_x_x" {
		t.Error("Unexpected string")
	}
}

func TestHookRelease(t *testing.T) {
	s := SimpleHook{L_INFO, bytes.NewBuffer(nil)}
	s2 := SimpleHook{L_INFO, bytes.NewBuffer(nil)}
	l := New(log.New(os.Stdout, "", log.LstdFlags), L_TRACE)
	l.Add(&s)
	l.Add(&s2)
	l.Trace("%s", func() string { return "x_x_x" })
	l.Release(&s)
	if l := len(l.hooks); l != 1 {
		t.Error("Failed", l)
	}
	l.Release(&s2)
	if l := len(l.hooks); l != 0 {
		t.Error("Failed", l)
	}
}
