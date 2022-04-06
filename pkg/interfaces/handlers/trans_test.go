package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/trans/pkg/domain"
	"github.mpi-internal.com/Yapo/trans/pkg/usecases"
)

type MockTransInteractor struct {
	mock.Mock
}

func (m *MockTransInteractor) ExecuteCommand(command domain.TransCommand) (domain.TransResponse, error) {
	ret := m.Called(command)
	return ret.Get(0).(domain.TransResponse), ret.Error(1)
}

func MakeMockInputTransGetter(input HandlerInput, response *goutils.Response) InputGetter {
	return func() (HandlerInput, *goutils.Response) {
		return input, response
	}
}

func TestTransHandlerInput(t *testing.T) {
	m := MockTransInteractor{}
	h := TransHandler{Interactor: &m}
	input := h.Input()
	var expected *TransHandlerInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
}

func TestTransHandlerExecuteOK(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{Command: "transinfo"}
	command := domain.TransCommand{
		Command: "transinfo",
		Params:  make([]domain.TransParams, 0),
	}
	response := domain.TransResponse{
		Status: usecases.TransOK,
		Params: make([]map[string]string, 1),
	}
	m.On("ExecuteCommand", command).Return(response, nil).Once()
	h := TransHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusOK,
		Body: TransRequestOutput{
			Status:   usecases.TransOK,
			Response: response.Params[0],
		},
	}

	getter := MakeMockInputTransGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestTransHandlerParseInput(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{
		Command: "get_account",
		Params:  make(map[string]interface{}),
	}
	input.Params["email"] = "user@test.com"
	command := domain.TransCommand{
		Command: "get_account",
		Params:  make([]domain.TransParams, 0),
	}
	param := domain.TransParams{
		Key:   "email",
		Value: "user@test.com",
	}
	command.Params = append(command.Params, param)

	response := domain.TransResponse{
		Status: usecases.TransOK,
		Params: make([]map[string]string, 1),
	}
	response.Params = []map[string]string{
		{
			"account_id": "1",
			"email":      "user@test.com",
			"is_company": "true",
		},
	}
	m.On("ExecuteCommand", command).Return(response, nil).Once()
	h := TransHandler{Interactor: &m}

	requestOutput := TransRequestOutput{
		Status: usecases.TransOK,
		Response: map[string]string{
			"account_id": "1",
			"email":      "user@test.com",
			"is_company": "true",
		},
	}
	expectedResponse := &goutils.Response{
		Code: http.StatusOK,
		Body: requestOutput,
	}

	getter := MakeMockInputTransGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestTransHandlerExecuteError(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{Command: "get_account"}
	command := domain.TransCommand{
		Command: "get_account",
		Params:  make([]domain.TransParams, 0),
	}
	response := domain.TransResponse{
		Status: usecases.TransError,
		Params: make([]map[string]string, 1),
	}
	m.On("ExecuteCommand", command).Return(response, nil).Once()
	h := TransHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: TransRequestOutput{
			Status:   usecases.TransError,
			Response: response.Params[0],
		},
	}

	getter := MakeMockInputTransGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestTransHandlerExecuteInternalError(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{Command: "get_account"}
	command := domain.TransCommand{
		Command: "get_account",
		Params:  make([]domain.TransParams, 0),
	}
	response := domain.TransResponse{
		Params: []map[string]string{{}},
	}
	m.On("ExecuteCommand", command).Return(response, errors.New("Error")).Once()
	h := TransHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusInternalServerError,
		Body: &goutils.GenericError{
			ErrorMessage: "Error",
		},
	}

	getter := MakeMockInputTransGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestTransHandlerInputError(t *testing.T) {
	m := MockTransInteractor{}
	h := TransHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: "Error",
	}

	getter := MakeMockInputTransGetter(nil, expectedResponse)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}
