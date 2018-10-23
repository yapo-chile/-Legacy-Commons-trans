package handlers

import (
	"github.com/Yapo/goutils"
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
	"net/http"
)

// FibonacciHandler implements the handler interface and responds to
// /fibonacci requests using an interactor. It's purpose is just to
// demonstrate Clean Architecture with a practical scenario
type FibonacciHandler struct {
	Interactor usecases.GetNthFibonacciUsecase
}

type fibonacciRequestInput struct {
	N int `json:"n"`
}

type fibonacciRequestOutput struct {
	Result domain.Fibonacci `json:"result"`
}

type fibonacciRequestError goutils.GenericError

// Input returns a fresh, empty instance of fibonacciRequestInput
func (h *FibonacciHandler) Input() HandlerInput {
	return &fibonacciRequestInput{}
}

// Execute carries on a /fibonacci request. Uses the given interactor to carry out
// the operation and get the desired value. Expected body format:
//	{
//		n: int - Number of fibonacci to retrieve (1 based)
//	}
// Expected response format:
//   { Result: int - Operation result }
// Expected error format:
//   { ErrorMessage: string - Error detail }
func (h *FibonacciHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}

	in := input.(*fibonacciRequestInput)
	f, err := h.Interactor.GetNth(in.N)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: fibonacciRequestError{
				err.Error(),
			},
		}
	}
	return &goutils.Response{
		Code: http.StatusOK,
		Body: fibonacciRequestOutput{
			Result: f,
		},
	}
}
