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
	Lvl
	*bytes.Buffer
}

func (s *SimpleHook) Id() string { return fmt.Sprintf("simple-hook-%d", s.N) }

func (s *SimpleHook) Match(l Lvl, msg, f string, a ...interface{}) bool { return l == s.Lvl }

func (s *SimpleHook) Action(l Lvl, msg, f string, a ...interface{}) { s.WriteString(msg + "\n") }

// Testing default functions
func TestDefault(t *testing.T) {
	Trace("trace msg")
	Debug("debug msg")
	Info("info msg")
	Error("error msg")
	Warn("warn msg")
	Panic("panic msg")
}

// Testing std console functions
func TestStd(t *testing.T) {
	c := Std()
	c.Trace("trace msg")
	c.Debug("debug msg")
	c.Info("info msg")
	c.Error("error msg")
	c.Warn("warn msg")
	c.Panic("panic msg")
}

func TestIgnored(t *testing.T) {
	SetDefaultCfg(Cfg{Lvl: LvlDebug})
	Trace("ignore msg")
	Debug("debug msg")
}

// Testing a hook with level error and a function argument
func TestHook(t *testing.T) {
	s := SimpleHook{0, LvlError, bytes.NewBuffer(nil)}
	l := New(Cfg{Lvl: LvlDebug, Color: true}, os.Stdout)
	l.Add(&s)
	l.Trace("mesage ignored by the logger")
	l.Error("message - %s", func() string { return "args" })
	l.Panic("ignored by the hook")
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

func TestClone(t *testing.T) {
	b := bytes.NewBuffer(nil)
	l := New(Cfg{Lvl: LvlInfo}, b)
	c := l.Clone("<prefix>")
	c.Debug("%s", "a")
	c.Warn("%s", "a")
	r := b.String()
	if exp := "WARN  <prefix> a\n"; r != exp {
		t.Errorf("Want %q, got %q", exp, r)
	}
	Std().Clone("prefix").Info("format")
}

func TestFormat(t *testing.T) {
	var cfgs = []Cfg{
		Cfg{File: FileHide, Date: DateHide},
		Cfg{File: FileShow, Date: DateHour},
		Cfg{File: FileFull, Date: DateFull},
	}
	for _, cfg := range cfgs {
		New(cfg, os.Stdout).Info("test a format")
	}
}

func TestPanic(t *testing.T) {
	var cfgs = []Cfg{
		Cfg{File: FileFmt(10)},
		Cfg{Date: DateFmt(10)},
	}
	for _, cfg := range cfgs {
		testPanic(t, cfg)
	}
}

func testPanic(t *testing.T, cfg Cfg) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	New(cfg, os.Stdout).Info("this will panic")
}
