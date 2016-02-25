#Console
[![GoDoc](https://godoc.org/gopkg.in/klaidliadon/console.v1?status.svg)](https://godoc.org/gopkg.in/klaidliadon/console.v1)
[![Build Status](https://travis-ci.org/klaidliadon/console.svg)](https://travis-ci.org/klaidliadon/console)
[![codecov.io](http://codecov.io/github/klaidliadon/console/coverage.svg?branch=master)](http://codecov.io/github/klaidliadon/console?branch=master)


Console package implements multi priority logger.

## Usage

It's recommended to use `gopkg.in` to ensure to use a stable version.

```go
import (
	"gopkg.in/klaidliadon/console.v1"
)
```


### Basic Usage

The standard console uses `os.Stdout`:

```go
c := console.Std()
c.Info("This is console")
```

### Custom Console

You can define your custom logger and use it:

```go
// Create a *console.Console
var custom = console.New(console.Cfg{
	Color: true, 
	Date: console.DateFull,
}, w)
custom.Trace("Ignored message %d", 1)
custom.Info("Message not ignored %d", 1)
```

## Features

### Runtime execution

You can use a `func() string` as argument. Instead
of the function the result string will be printed.
If the message is ignored the function will not be executed.

```go
var a = Tree{} // With a very expensive String method

l := console.New(console.Cfg{Lvl: console.LvlDebug})
// Tree.String is executed
l.Info("Method result: %s", a)
// Tree.String is ignored
l.Trace("Method result: %s", a)
```

### Hooks

An hook is interface used to capture certain conditions and execute an action.

Here's an example:

```go
type Mailer struct {
	lvl  console.Lvl
	Addr string
	Auth smtp.Auth
	From string
	To   []string
}

func (m Mailer) Match(l console.Lvl, format string, args ...interface{}) bool {
	return l >= m.lvl
}

func (m Mailer) Action(l console.Lvl, msg string){
	smtp.SendMail(m.Addr, m.Auth, m.From, m.To, fmt.Sprintf("[%s] from MailHook: %s", l, msg)
}
```

This hook captures the messages from a certain level and sends an email with the message content.

## Help

For a complete reference read the [docs](http://godoc.org/gopkg.in/klaidliadon/console.v1 "Godoc").
