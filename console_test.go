package console

import "testing"

func TestConsole(t *testing.T) {
	l := Std(L_TRACE)
	var logFunc = func(s string) {
		l.Output(1, L_WARN, s)
	}
	l.Output(0, 0, "a")
	l.Info("some text %q", "a")
	logFunc("text")
	l.Clone("[prefix]").Info("some text %q", "a")
}
