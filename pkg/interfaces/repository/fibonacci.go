package repository

import (
	"fmt"
	"sync"

	"github.schibsted.io/Yapo/goms/pkg/domain"
)

// mapFibonacciRepository is an implementation of domain.FibonacciRepository
// that stores Fibonacci data on a map. It keeps the lastest known pair on a
// separate array to speed up retrieval. The type is intentionally private.
// The correct way to instantiate this type is with NewMapFibonacciRepository.
// This ensures that the required initialization is performed every time.
type mapFibonacciRepository struct {
	storage map[int]domain.Fibonacci
	latest  []int
	mutex   sync.RWMutex
}

// NewMapFibonacciRepository instantiates a fresh mapFibonacciRepository,
// performs the initialization and returns it as a domain.FibonacciRepository.
// The return type prevents others to directly access data members.
func NewMapFibonacciRepository() domain.FibonacciRepository {
	var r mapFibonacciRepository
	r.storage = map[int]domain.Fibonacci{
		1: 1,
		2: 1,
	}
	r.latest = []int{1, 2}
	return &r
}

// Get returns the nth (1 based) Fibonacci should this instance know it.
// Otherwise, will return -1 and error message.
func (r *mapFibonacciRepository) Get(nth int) (domain.Fibonacci, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	f, found := r.storage[nth]
	if !found {
		return -1, fmt.Errorf("Don't know the %dth Fibonacci, do you?", nth)
	}
	return f, nil
}

// Save sets the nth Fibonacci to x should the last known pair of values end
// at nth-1. Otherwise returns an error message.
func (r *mapFibonacciRepository) Save(nth int, x domain.Fibonacci) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if nth != r.latest[1]+1 {
		return fmt.Errorf("How do you know the %dth Fibonacci number?", nth)
	}
	r.storage[nth] = x
	r.latest[0]++
	r.latest[1]++
	return nil
}

// LatestPair retrieves the latest pair (and indexes) of known Fibonacci Numbers
func (r *mapFibonacciRepository) LatestPair() domain.FibonacciPair {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return domain.FibonacciPair{
		IA: r.latest[0],
		IB: r.latest[1],
		A:  r.storage[r.latest[0]],
		B:  r.storage[r.latest[1]],
	}
}
