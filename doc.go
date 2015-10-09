/*

Console package implements multi priority logger.

The standard log uses os.Stdout:

	c := console.Std()
	c.Info("This is console")

You can define your custom logger and use it:

	var custom = console.New(console.Cfg{ Color: true, Date: console.DateFull}, w)
	custom.Trace("Ignored message %d", 1)
	custom.Info("Message not ignored %d", 1)

You can use a `func() string` as argument. Instead
of the function the result string will be printed.
If the message is ignored the function will not be executed.

	// Very expensive struct method
	var a func() string = myObject.createTreeString
	// Very expensive interface method
	var b func() string = myInterface.createTreeString

	l := logger.New(console.Cfg{Lvl: console.LvlDebug})
	// func a is executed
	l.Info("Method result: %s", a)
	// func b is ignored
	l.Trace("Method result: %s", b)

A main feature of the package is Hooks: a hook is interface used to capture certain conditions and execute an action.
Here's an example:

	type MailHook struct {
		lvl  console.LogLvl
		Addr string
		Auth Auth
		From string
		To   []string
	}

	func (m *MailHook) Match(l console.LogLvl, format string, args ...interface{}) bool {
		return l >= m.lvl
	}

	func (m *MailHook) Action(l LogLvl, msg string){
		smtp.SendMail(m.Addr, m.Auth, m.From, m.To, fmt.Sprintf("[%s] from MailHook: %s\n\n%s", l, fileline, msg)
	}


This hook captures the messages from a certain level and sends an email with the message content.

*/
package console
