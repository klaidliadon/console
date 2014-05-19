package console

import (
	"bytes"
	"testing"
)

func TestConsole(t *testing.T) {
	l := Std(L_TRACE)
	var logFunc = func(s string) {
		l.Output(1, L_WARN, s)
	}
	s := SimpleHook{L_ERROR, bytes.NewBuffer(nil)}
	l.Add(&s)
	l.Output(0, 0, "a")
	l.Info("some text %q", "a")
	logFunc("text")
	l.Clone("[prefix]").Info("some text %q", "a")
	l.Error("message - %s", "args")
	t.Log(s.String())
}
