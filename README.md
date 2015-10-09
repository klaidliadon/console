#Console
[![GoDoc](https://godoc.org/gopkg.in/klaidliadon/console.v1?status.svg)](https://godoc.org/gopkg.in/klaidliadon/console.v1) 
[![codecov.io](http://codecov.io/github/klaidliadon/console/coverage.svg?branch=master)](http://codecov.io/github/klaidliadon/console?branch=master)


Console package implements multi priority logger.

## Usage

### Basic Usage

The standard log uses `os.Stdout`:

	c := console.Std()
	c.Info("This is console")

### Custom Logger

You can define your custom logger and use it:

	// Create a *log.Logger
	var custom = console.New(console.Cfg{
		Color: true, 
		Date: console.DateFull,
	}, w)
	custom.Trace("Ignored message %d", 1)
	custom.Info("Message not ignored %d", 1)

## Features

### Runtime execution

You can use a `func() string` as argument. Instead
of the function the result string will be printed.
If the message is ignored the function will not be executed.

	// Very expensive struct method
	var a func() string = myObject.createTreeString
	// Very expensive interface method
	var b func() string = myInterface.createTreeString

	l := console.New(console.Cfg{Lvl: console.LvlDebug})
	// func a is executed
	l.Info("Method result: %s", a)
	// func b is ignored
	l.Trace("Method result: %s", b)

### Hooks

An hook is interface used to capture certain conditions and execute an action.

Here's an example:

	type MailHook struct {
		lvl  console.LogLevel
		Addr string
		Auth Auth
		From string
		To   []string
	}

	func (m *MailHook) Match(l console.LogLevel, format string, args ...interface{}) bool {
		return l >= m.lvl
	}

	func (m *MailHook) Action(l LogLevel, msg string){
		smtp.SendMail(m.Addr, m.Auth, m.From, m.To, fmt.Sprintf("[%s] from MailHook: %s\n\n%s", l, fileline, msg)
	}


This hook captures the messages from a certain level and sends an email with the message content.

## Help

For a complete reference read the [docs](http://godoc.org/gopkg.in/klaidliadon/console.v1 "Godoc").
