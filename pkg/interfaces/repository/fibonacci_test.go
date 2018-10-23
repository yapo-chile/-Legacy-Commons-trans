package repository

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.schibsted.io/Yapo/goms/pkg/domain"
)

func TestFibonacciRepositoryCreation(t *testing.T) {
	r := NewMapFibonacciRepository()

	x, err := r.Get(1)
	assert.Equal(t, domain.Fibonacci(1), x)
	assert.NoError(t, err)

	x, err = r.Get(2)
	assert.Equal(t, domain.Fibonacci(1), x)
	assert.NoError(t, err)

	x, err = r.Get(3)
	assert.Equal(t, domain.Fibonacci(-1), x)
	assert.Error(t, err)

	p := r.LatestPair()
	expected := domain.FibonacciPair{
		IA: 1, A: domain.Fibonacci(1),
		IB: 2, B: domain.Fibonacci(1),
	}

	assert.Equal(t, expected, p)
}

func TestFibonacciRepositorySaveWildGuess(t *testing.T) {
	r := NewMapFibonacciRepository()

	err := r.Save(5, 32)
	assert.Error(t, err)
}

func TestFibonacciRepositorySaveReplace(t *testing.T) {
	r := NewMapFibonacciRepository()

	err := r.Save(2, 42)
	assert.Error(t, err)
}

func TestFibonacciRepositorySaveNext(t *testing.T) {
	r := NewMapFibonacciRepository()

	err := r.Save(3, 2)
	assert.NoError(t, err)

	p := r.LatestPair()
	expected := domain.FibonacciPair{
		IA: 2, A: domain.Fibonacci(1),
		IB: 3, B: domain.Fibonacci(2),
	}

	assert.Equal(t, expected, p)
}

func TestFibonacciRepositoryConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	r := NewMapFibonacciRepository()
	f := func(r domain.FibonacciRepository) {
		defer wg.Done()
		a, _ := r.Get(1)
		b, _ := r.Get(2)
		for i := 3; i < 1000; i++ {
			r.Save(i, a+b)
			p := r.LatestPair()
			a = p.A
			b, _ = r.Get(i)
		}
	}
	wg.Add(2)
	go f(r)
	go f(r)
	wg.Wait()
}
