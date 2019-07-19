package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/trans/pkg/domain"
	"github.mpi-internal.com/Yapo/trans/pkg/usecases"
)

type MockTransHandler struct {
	mock.Mock
}

func (m *MockTransHandler) SendCommand(command string, params map[string]string) (map[string]string, error) {
	ret := m.Called(command, params)
	return ret.Get(0).(map[string]string), ret.Error(1)
}

type MockTransFactory struct {
	mock.Mock
}

func (m *MockTransFactory) MakeTransHandler() TransHandler {
	ret := m.Called()
	return ret.Get(0).(TransHandler)
}

func TestNewTransRepo(t *testing.T) {
	factory := MockTransFactory{}
	repo := NewTransRepo(&factory)

	expectedRepo := &TransRepo{
		transFactory: &factory,
	}
	assert.Equal(t, expectedRepo, repo)
	factory.AssertExpectations(t)
}

func TestExecuteError(t *testing.T) {
	cmd := "command1"
	params := make(map[string]string)
	expectedErr := errors.New("trans error")
	responseParams := make(map[string]string)

	command := domain.TransCommand{
		Command: cmd,
		Params:  make([]domain.TransParams, 0),
	}

	handler := MockTransHandler{}
	handler.On("SendCommand", cmd, params).Return(responseParams, expectedErr).Once()

	factory := MockTransFactory{}
	factory.On("MakeTransHandler").Return(&handler)

	repo := NewTransRepo(&factory)

	response, err := repo.Execute(command)
	expectedResponse := domain.TransResponse{
		Params: make(map[string]string),
	}
	expectedResponse.Params["error"] = "trans error"
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedResponse, response)
	factory.AssertExpectations(t)
	handler.AssertExpectations(t)
}

func TestExecuteOK(t *testing.T) {
	cmd := "command1"
	params := make(map[string]string)
	params["param 1"] = "value 1"
	params["param 2"] = "value 2"

	responseParams := make(map[string]string)
	responseParams["status"] = usecases.TransOK
	responseParams["response 1"] = "response 1"
	command := domain.TransCommand{
		Command: cmd,
		Params:  make([]domain.TransParams, 0),
	}
	for key, val := range params {
		command.Params = append(
			command.Params,
			domain.TransParams{
				Key:   key,
				Value: val,
			},
		)
	}

	handler := MockTransHandler{}
	handler.On("SendCommand", cmd, params).Return(responseParams, nil).Once()

	factory := MockTransFactory{}
	factory.On("MakeTransHandler").Return(&handler).Once()

	repo := NewTransRepo(&factory)

	response, err := repo.Execute(command)
	expectedResponse := domain.TransResponse{
		Status: usecases.TransOK,
		Params: make(map[string]string),
	}
	expectedResponse.Params["response 1"] = "response 1"
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
	factory.AssertExpectations(t)
	handler.AssertExpectations(t)
}

func TestExecuteOKNumbers(t *testing.T) {
	cmd := "command1"
	params := make(map[string]string)
	params["param 1"] = "1980"

	responseParams := make(map[string]string)
	responseParams["status"] = usecases.TransOK
	responseParams["response 1"] = "response 1"
	command := domain.TransCommand{
		Command: cmd,
		Params:  make([]domain.TransParams, 0),
	}
	command.Params = append(
		command.Params,
		domain.TransParams{
			Key:   "param 1",
			Value: 1980,
		},
	)

	handler := MockTransHandler{}
	handler.On("SendCommand", cmd, params).Return(responseParams, nil).Once()

	factory := MockTransFactory{}
	factory.On("MakeTransHandler").Return(&handler).Once()

	repo := NewTransRepo(&factory)

	response, err := repo.Execute(command)
	expectedResponse := domain.TransResponse{
		Status: usecases.TransOK,
		Params: make(map[string]string),
	}
	expectedResponse.Params["response 1"] = "response 1"
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
	factory.AssertExpectations(t)
	handler.AssertExpectations(t)
}
