package services

import (
	"reflect"
	"strconv"

	"github.mpi-internal.com/Yapo/trans/pkg/domain"
)

// TransHandler is an interface to use Trans functions
type TransHandler interface {
	SendCommand(string, []domain.TransParams) (domain.TransResponse, error)
}

// TransFactory is an interface that abstracts the Factory Pattern for creating TransHandler objects
type TransFactory interface {
	MakeTransHandler() TransHandler
}

// TransRepo struct definition
type TransRepo struct {
	transFactory TransFactory
}

// NewTransRepo instance TransRepo and set handler
func NewTransRepo(transFactory TransFactory) *TransRepo {
	return &TransRepo{
		transFactory: transFactory,
	}
}

// Execute executes the specified trans command
func (repo *TransRepo) Execute(command domain.TransCommand) (domain.TransResponse, error) {
	return repo.transaction(command.Command, command.Params)
}

func (repo *TransRepo) transaction(method string, transParams []domain.TransParams) (domain.TransResponse, error) {
	trans := repo.transFactory.MakeTransHandler()
	for _, transParam := range transParams {
		if reflect.TypeOf(transParam.Value).Kind() == reflect.Int {
			transParam.Value = strconv.Itoa(transParam.Value.(int))
		}
	}
	return trans.SendCommand(method, transParams)
}
