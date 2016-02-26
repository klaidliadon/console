package console

import (
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// Date format type
type DateFmt int

func (d DateFmt) fmt() func(t time.Time) string {
	switch d {
	case DateHide:
		return nil
	case DateHour:
		return func(t time.Time) string { return t.Format("15:04:05") }
	case DateFull:
		return func(t time.Time) string { return t.Format("2006/01/02 15:04:05") }
	default:
		panic("Invalid DateFmt")
	}
}

// All the date format configurations
const (
	DateHide = DateFmt(iota)
	DateHour
	DateFull
)

// Filename format type
type FileFmt int

func (f FileFmt) fmt() func(n string) string {
	switch f {
	case FileHide:
		return nil
	case FileShow:
		return func(n string) string {
			d, f := filepath.Split(n)
			return filepath.Base(d) + string(filepath.Separator) + f
		}
	case FileFull:
		return func(n string) string { return n }
	default:
		panic("Invalid FileFmt")
	}
}

// All the file path configurations
const (
	FileHide = FileFmt(iota)
	FileShow
	FileFull
)

// Logging level
type Lvl int

// List of priorities
const (
	LvlTrace = Lvl(iota)
	LvlDebug
	LvlInfo
	LvlWarn
	LvlError
	LvlPanic
)

type desc struct {
	label string
	Color func(...interface{}) string
}

func (l desc) Label(c bool) string {
	if c {
		return l.Color(l.label)
	}
	return l.label
}

func newDesc(s string, c *color.Color) desc {
	return desc{s, c.SprintFunc()}
}

var levels = map[Lvl]desc{
	LvlTrace: newDesc("TRACE", color.New(color.FgHiBlue)),
	LvlDebug: newDesc("DEBUG", color.New(color.FgHiCyan)),
	LvlInfo:  newDesc("INFO ", color.New(color.FgHiGreen)),
	LvlWarn:  newDesc("WARN ", color.New(color.FgHiYellow)),
	LvlError: newDesc("ERROR", color.New(color.FgHiRed)),
	LvlPanic: newDesc("PANIC", color.New(color.FgHiMagenta)),
}

// Holds the configuration of a Console
type Cfg struct {
	Date   DateFmt
	File   FileFmt
	Lvl    Lvl
	Color  bool
	prefix string
	fmt    struct {
		date func(time.Time) string
		file func(string) string
	}
}

func (c *Cfg) validate() {
	c.fmt.date = c.Date.fmt()
	c.fmt.file = c.File.fmt()
}
