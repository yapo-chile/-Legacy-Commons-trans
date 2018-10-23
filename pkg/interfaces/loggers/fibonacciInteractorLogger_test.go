package loggers

import (
	"testing"
)

// There are no return values to assert on, as logger only cause side effects
// to communicate with the outside world. These tests only ensure that the
// loggers don't panic

func TestFibonacciInteractorDefaultLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeFibonacciInteractorLogger(m)
	l.LogBadInput(42)
	l.LogRepositoryError(5, 42, nil)
}
