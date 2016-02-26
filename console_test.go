package console

import (
	"bytes"
	"os"

	. "gopkg.in/check.v1"
)

var _ = Suite(ConsoleSuite{})

type ConsoleSuite struct{}

func (s ConsoleSuite) TestDefault(c *C) {
	SetDefaultCfg(Cfg{})
	Trace("trace msg")
	Debug("debug msg")
	Info("info msg")
	Error("error msg")
	Warn("warn msg")
	Panic("panic msg")
}

func (s ConsoleSuite) TestIgnored(c *C) {
	b := bytes.NewBuffer(nil)
	console := New(Cfg{Lvl: LvlDebug}, b)
	console.Trace("ignore msg")
	c.Assert(b.Len(), Equals, 0)
	console.Debug("debug msg")
	c.Assert(b.Len(), Not(Equals), 0)
}

func (s ConsoleSuite) TestHook(c *C) {
	h := SimpleHook{0, LvlError, bytes.NewBuffer(nil)}
	l := New(Cfg{Lvl: LvlDebug, Color: true}, os.Stdout)
	l.Add(&h)
	l.Trace("mesage ignored by the logger")
	l.Error("message - %s", func() string { return "args" })
	l.Panic("ignored by the hook")
	c.Assert(h.String(), Equals, "message - args")
}

func (s ConsoleSuite) TestFunction(c *C) {
	h := SimpleHook{0, LvlInfo, bytes.NewBuffer(nil)}
	v := "text\n"
	l := Std()
	l.Add(&h)
	l.Info("%s", func() string { return v })
	c.Assert(h.String(), Equals, v)
}

func (s ConsoleSuite) TestHookRelease(c *C) {
	h1 := SimpleHook{1, LvlInfo, bytes.NewBuffer(nil)}
	h2 := SimpleHook{2, LvlInfo, bytes.NewBuffer(nil)}
	l := New(Cfg{Lvl: LvlTrace, Color: true}, os.Stdout)
	l.Add(&h1)
	l.Add(&h2)
	c.Assert(l.hooks, HasLen, 2)
	l.Release(&h1)
	c.Assert(l.hooks, HasLen, 1)
	l.Release(&h2)
	c.Assert(l.hooks, HasLen, 0)
}

func (s ConsoleSuite) TestClone(c *C) {
	b := bytes.NewBuffer(nil)
	l := New(Cfg{Lvl: LvlInfo}, b)
	clone := l.Clone("<prefix>")
	c.Assert(l.cfg.prefix, Not(Equals), clone.cfg.prefix)
	clone.Debug("%s", "a")
	clone.Warn("%s", "a")
	c.Assert(b.String(), Equals, "WARN  <prefix> a\n")
	l.Debug("%s", "a")
	l.Warn("%s", "a")
	Std().Clone("prefix").Info("format")
}

func (s ConsoleSuite) TestFormat(c *C) {
	var cfgs = []Cfg{
		Cfg{Color: true, File: FileHide, Date: DateHide},
		Cfg{Color: true, File: FileShow, Date: DateHour},
		Cfg{Color: true, File: FileFull, Date: DateFull},
	}
	for _, cfg := range cfgs {
		New(cfg, os.Stdout).Info("test a format")
	}
}

func (s ConsoleSuite) TestPanic(c *C) {
	var cfgs = []struct {
		Cfg
		Panic interface{}
	}{
		{Cfg{File: FileFmt(10)}, "Invalid FileFmt"},
		{Cfg{Date: DateFmt(10)}, "Invalid DateFmt"},
	}
	for _, cfg := range cfgs {
		c.Assert(func() {
			New(cfg.Cfg, os.Stdout).Info("this will panic")
		}, Panics, cfg.Panic)
	}
}

func (s ConsoleSuite) BenchmarkPlain(c *C) {
	l := New(Cfg{}, null{})
	for n := 0; n < c.N; n++ {
		l.Info("msg %v", n)
	}
}

func (s ConsoleSuite) BenchmarkStd(c *C) {
	l := New(Defaults, null{})
	for n := 0; n < c.N; n++ {
		l.Info("msg %v", n)
	}
}

func (s ConsoleSuite) BenchmarkStdNoColor(c *C) {
	var cfg = Defaults
	cfg.Color = false
	l := New(cfg, null{})
	for n := 0; n < c.N; n++ {
		l.Info("msg %v", n)
	}
}
