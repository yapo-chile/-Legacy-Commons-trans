package loggers

import (
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
)

type fibonacciInteractorDefaultLogger struct {
	logger Logger
}

func (l *fibonacciInteractorDefaultLogger) LogBadInput(n int) {
	l.logger.Debug("GetNth doesn't like N < 1. Input: %d", n)
}

func (l *fibonacciInteractorDefaultLogger) LogRepositoryError(i int, x domain.Fibonacci, err error) {
	l.logger.Error("Repository refused to save (%d, %d): %s", i, x, err)
}

// MakeFibonacciInteractorLogger sets up a FibonacciInteractorLogger instrumented
// via the provided logger
func MakeFibonacciInteractorLogger(logger Logger) usecases.FibonacciInteractorLogger {
	return &fibonacciInteractorDefaultLogger{
		logger: logger,
	}
}
