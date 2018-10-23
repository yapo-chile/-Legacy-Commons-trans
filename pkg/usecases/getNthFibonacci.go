package usecases

import (
	"fmt"
	"github.schibsted.io/Yapo/goms/pkg/domain"
)

// GetNthFibonacciUsecase states:
// As a User, I would like to know which the Nth Fibonacci Number is.
// GetNth should return that number to me, or an appropiate error if not possible.
type GetNthFibonacciUsecase interface {
	GetNth(n int) (domain.Fibonacci, error)
}

// FibonacciInteractorLogger defines all the events a FibonacciInteractor may
// need/like to report as they happen
type FibonacciInteractorLogger interface {
	LogBadInput(int)
	LogRepositoryError(int, domain.Fibonacci, error)
}

// FibonacciInteractor implements GetNthFibonacciUsecase by using Repository
// to store new Fibonacci as required and to retrieve the final answer.
type FibonacciInteractor struct {
	Logger     FibonacciInteractorLogger
	Repository domain.FibonacciRepository
}

// GetNth finds the nth Fibonacci Number by recursively generating one more
// from the last know pair. The running time is O(n).
func (interactor *FibonacciInteractor) GetNth(n int) (domain.Fibonacci, error) {
	// Ensure correct input
	if n <= 0 {
		interactor.Logger.LogBadInput(n)
		return -1, fmt.Errorf("There's no such thing as %dth Fibonacci", n)
	}
	// Check if the repository already knows it
	x, err := interactor.Repository.Get(n)
	if err == nil {
		return x, nil
	}
	// Retrieve the latest pair
	latest := interactor.Repository.LatestPair()
	i, x := latest.Next()
	err = interactor.Repository.Save(i, x)
	if err != nil {
		// Report the error
		interactor.Logger.LogRepositoryError(i, x, err)
		return -1, err
	}
	// One step closer. Keep trying
	return interactor.GetNth(n)
}
