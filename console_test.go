package console

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

// A simple hook that copies in a buffer messages from a certain level.
type SimpleHook struct {
	N int
	LogLvl
	*bytes.Buffer
}

func (s *SimpleHook) Id() string { return fmt.Sprintf("simple-hook-%d", s.N) }

func (s *SimpleHook) Match(l LogLvl, f string, a ...interface{}) bool { return l == s.LogLvl }

func (s *SimpleHook) Action(l LogLvl, msg string) { s.WriteString(msg + "\n") }

// Testing default functions
func TestDefault(t *testing.T) {
	Trace("trace msg")
	Debug("debug msg")
	Info("info msg")
	Error("error msg")
	Warn("warn msg")
	Panic("panic msg")
}

// Testing a hook with level error and a function argument
func TestHook(t *testing.T) {
	s := SimpleHook{0, LvlError, bytes.NewBuffer(nil)}
	l := New(Cfg{Lvl: LvlTrace, Color: true}, os.Stdout)
	l.Add(&s)
	l.Error("message - %s", func() string { return "args" })
	l.Panic("ignored")
	if s.String() != "message - args\n" {
		t.Error("Unexpected string", s)
	}
}

func TestFunction(t *testing.T) {
	s := SimpleHook{0, LvlInfo, bytes.NewBuffer(nil)}
	l := Std()
	l.Add(&s)
	l.Info("%s", func() string { return "x_x_x" })
	if s.String() != "x_x_x\n" {
		t.Error("Unexpected string")
	}
}

func TestHookRelease(t *testing.T) {
	s1 := SimpleHook{1, LvlInfo, bytes.NewBuffer(nil)}
	s2 := SimpleHook{2, LvlInfo, bytes.NewBuffer(nil)}
	l := New(Cfg{Lvl: LvlTrace, Color: true}, os.Stdout)
	l.Add(&s1)
	l.Add(&s2)
	l.Trace("%s", func() string { return "x_x_x" })
	l.Release(&s1)
	if l := len(l.hooks); l != 1 {
		t.Error("Failed", l)
	}
	l.Release(&s2)
	if l := len(l.hooks); l != 0 {
		t.Error("Failed", l)
	}
}
