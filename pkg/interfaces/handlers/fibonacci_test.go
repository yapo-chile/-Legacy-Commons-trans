package handlers

import (
	"errors"
	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"net/http"
	"testing"
)

type MockFibonacciInteractor struct {
	mock.Mock
}

func (m *MockFibonacciInteractor) GetNth(n int) (domain.Fibonacci, error) {
	args := m.Called(n)
	return args.Get(0).(domain.Fibonacci), args.Error(1)
}

func MakeMockInputGetter(input HandlerInput, response *goutils.Response) InputGetter {
	return func() (HandlerInput, *goutils.Response) {
		return input, response
	}
}

func TestFibonacciHandlerInput(t *testing.T) {
	m := MockFibonacciInteractor{}
	h := FibonacciHandler{Interactor: &m}
	input := h.Input()
	var expected *fibonacciRequestInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
}

func TestFibonacciHandlerExecuteOK(t *testing.T) {
	m := MockFibonacciInteractor{}
	m.On("GetNth", 5).Return(domain.Fibonacci(5), nil).Once()
	h := FibonacciHandler{Interactor: &m}

	input := fibonacciRequestInput{N: 5}
	expectedResponse := &goutils.Response{
		Code: http.StatusOK,
		Body: fibonacciRequestOutput{5},
	}

	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestFibonacciHandlerExecuteError(t *testing.T) {
	m := MockFibonacciInteractor{}
	m.On("GetNth", -1).Return(domain.Fibonacci(0), errors.New("kaboom")).Once()
	h := FibonacciHandler{Interactor: &m}

	input := fibonacciRequestInput{N: -1}
	expectedResponse := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: fibonacciRequestError{"kaboom"},
	}

	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestFibonacciHandlerInputError(t *testing.T) {
	m := MockFibonacciInteractor{}
	h := FibonacciHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: fibonacciRequestError{"kaboom"},
	}

	getter := MakeMockInputGetter(nil, expectedResponse)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}
