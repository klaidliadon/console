package console

import "testing"

func TestConsole(t *testing.T) {
	l := Std(L_TRACE)
	var logFunc = func(s string) {
		l.Output(0, L_WARN, s)
	}
	l.Info("some text %q", "a")
	logFunc("text")
	l.Clone("[prefix]").Info("some text %q", "a")
}
