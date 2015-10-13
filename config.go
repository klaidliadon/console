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
type LogLvl int

// List of priorities
const (
	LvlTrace = LogLvl(iota)
	LvlDebug
	LvlInfo
	LvlWarn
	LvlError
	LvlPanic
)

type desc struct {
	Label string
	Color *color.Color
}

func (l desc) GetLabel(c bool) string {
	if c {
		return l.Color.SprintFunc()(l.Label)
	}
	return l.Label
}

var listOfLvls = map[LogLvl]desc{
	LvlTrace: desc{"TRACE", color.New(color.FgBlue)},
	LvlDebug: desc{"DEBUG", color.New(color.FgCyan)},
	LvlInfo:  desc{"INFO ", color.New(color.FgGreen)},
	LvlWarn:  desc{"WARN ", color.New(color.FgYellow)},
	LvlError: desc{"ERROR", color.New(color.FgRed)},
	LvlPanic: desc{"PANIC", color.New(color.FgMagenta)},
}

// Holds the configuration of a Console
type Cfg struct {
	Date   DateFmt
	File   FileFmt
	Lvl    LogLvl
	Color  bool
	prefix string
}
