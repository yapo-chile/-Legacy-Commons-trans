package usecases

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.schibsted.io/Yapo/goms/pkg/domain"
)

type MockFibonacciRepository struct {
	mock.Mock
}

func (m *MockFibonacciRepository) Get(nth int) (domain.Fibonacci, error) {
	ret := m.Called(nth)
	return ret.Get(0).(domain.Fibonacci), ret.Error(1)
}
func (m *MockFibonacciRepository) Save(nth int, x domain.Fibonacci) error {
	ret := m.Called(nth, x)
	return ret.Error(0)
}
func (m *MockFibonacciRepository) LatestPair() domain.FibonacciPair {
	ret := m.Called()
	return ret.Get(0).(domain.FibonacciPair)
}

type MockFibonacciLogger struct {
	mock.Mock
}

func (m *MockFibonacciLogger) LogBadInput(x int) {
	m.Called(x)
}
func (m *MockFibonacciLogger) LogRepositoryError(i int, x domain.Fibonacci, err error) {
	m.Called(i, x, err)
}

func TestFibonacciInteractorGetNthNegative(t *testing.T) {
	l := &MockFibonacciLogger{}
	m := &MockFibonacciRepository{}
	i := FibonacciInteractor{
		Logger:     l,
		Repository: m,
	}

	l.On("LogBadInput", -1)

	_, err := i.GetNth(-1)
	assert.Error(t, err)
	m.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestFibonacciInteractorGetNthKnown(t *testing.T) {
	l := &MockFibonacciLogger{}
	m := &MockFibonacciRepository{}
	m.On("Get", 1).Return(domain.Fibonacci(1), nil)

	i := FibonacciInteractor{
		Logger:     l,
		Repository: m,
	}

	x, err := i.GetNth(1)
	assert.Equal(t, domain.Fibonacci(1), x)
	assert.NoError(t, err)
	m.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestFibonacciInteractorGetNthUnknown(t *testing.T) {
	l := &MockFibonacciLogger{}
	m := &MockFibonacciRepository{}
	m.On("Get", 4).Return(domain.Fibonacci(-1), errors.New("Some error")).Once()
	m.On("LatestPair").Return(domain.FibonacciPair{IA: 1, A: domain.Fibonacci(1), IB: 2, B: domain.Fibonacci(1)}).Once()
	m.On("Save", 3, domain.Fibonacci(2)).Return(nil)

	m.On("Get", 4).Return(domain.Fibonacci(-1), errors.New("Some error")).Once()
	m.On("LatestPair").Return(domain.FibonacciPair{IA: 2, A: domain.Fibonacci(1), IB: 3, B: domain.Fibonacci(2)}).Once()
	m.On("Save", 4, domain.Fibonacci(3)).Return(nil)

	m.On("Get", 4).Return(domain.Fibonacci(3), nil)

	i := FibonacciInteractor{
		Logger:     l,
		Repository: m,
	}

	x, err := i.GetNth(4)
	assert.Equal(t, domain.Fibonacci(3), x)
	assert.NoError(t, err)
	m.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestFibonacciInteractorGetNthSaveError(t *testing.T) {
	l := &MockFibonacciLogger{}
	m := &MockFibonacciRepository{}
	m.On("Get", 4).Return(domain.Fibonacci(-1), errors.New("Some error")).Once()
	m.On("LatestPair").Return(domain.FibonacciPair{IA: 1, A: domain.Fibonacci(1), IB: 2, B: domain.Fibonacci(1)}).Once()
	m.On("Save", 3, domain.Fibonacci(2)).Return(errors.New("Weird error")).Once()

	l.On("LogRepositoryError", 3, domain.Fibonacci(2), errors.New("Weird error"))
	i := FibonacciInteractor{
		Logger:     l,
		Repository: m,
	}

	_, err := i.GetNth(4)
	assert.Error(t, err)
	m.AssertExpectations(t)
	l.AssertExpectations(t)
}
