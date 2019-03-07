package usecases

import (
	"fmt"
	"strings"

	"github.schibsted.io/Yapo/trans/pkg/domain"
)

// ExecuteTransUsecase states:
// As a User, I would like to execute my TransCommand on a Trans server and get the corresponding response
// ExecuteTrans should return a response, or an appropiate error if there was a problem.
type ExecuteTransUsecase interface {
	ExecuteCommand(command domain.TransCommand) (domain.TransResponse, error)
}

// TransInteractorLogger defines all the events a TransInteractor may
// need/like to report as they happen
type TransInteractorLogger interface {
	LogBadInput(domain.TransCommand)
	LogRepositoryError(domain.TransCommand, error)
}

// TransInteractor implements ExecuteTransUsecase by using Repository
// to execute the Trans and to retrieve the response.
type TransInteractor struct {
	Logger     TransInteractorLogger
	Repository domain.TransRepository
}

// ExecuteCommand executes the given TransCommand and returns the corresponding TransResponse.
func (interactor TransInteractor) ExecuteCommand(
	command domain.TransCommand,
) (domain.TransResponse, error) {
	response := domain.TransResponse{
		Status: "TRANS_ERROR",
		Params: make(map[string]string),
	}
	// Ensure correct input
	if command.Command == "" {
		interactor.Logger.LogBadInput(command)
		return response, fmt.Errorf("invalid command %+v", command)
	}

	// Execute the command and retrieve the response
	response, err := interactor.Repository.Execute(command)
	if err != nil {
		// Report the error
		interactor.Logger.LogRepositoryError(command, err)
		if transErr, ok := response.Params["error"]; ok {
			err = fmt.Errorf(transErr)
		} else {
			err = fmt.Errorf("error during execution")
		}
	}
	// if the command sent doesn´t exists in the server
	if response.Status == "TRANS_ERROR_NO_SUCH_COMMAND:Err no such command" {
		err = fmt.Errorf("error command doesn't exists")
		response.Status = "TRANS_ERROR"
		response.Params["error"] = err.Error()
	}
	// if the error is a database error
	if strings.Contains(response.Status, "TRANS_DATABASE_ERROR") {
		//get the specific error message from the status response
		err = fmt.Errorf(strings.Replace(response.Status, "TRANS_DATABASE_ERROR:", "", 1))
		interactor.Logger.LogRepositoryError(command, err)
		response.Status = "TRANS_DATABASE_ERROR"
		response.Params["error"] = err.Error()
	}

	return response, err
}
