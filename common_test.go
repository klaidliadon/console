package console

import (
	"bytes"
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type SimpleHook struct {
	N int
	Lvl
	*bytes.Buffer
}

func (s *SimpleHook) Id() string                                      { return fmt.Sprintf("simple-hook-%d", s.N) }
func (s *SimpleHook) Match(l Lvl, _, _ string, _ ...interface{}) bool { return l == s.Lvl }
func (s *SimpleHook) Action(_ Lvl, msg, _ string, _ ...interface{})   { s.WriteString(msg) }

type null struct{}

func (null) Write(p []byte) (n int, err error)       { return 0, nil }
func (null) WriteString(s string) (n int, err error) { return 0, nil }
