package handlers

import (
	"net/http"

	"github.com/Yapo/goutils"
	"github.mpi-internal.com/Yapo/trans/pkg/domain"
	"github.mpi-internal.com/Yapo/trans/pkg/usecases"
)

// TransHandler implements the handler interface and responds to /execute
// requests with a message. Expected response format:
// { status: string, response: json }
type TransHandler struct {
	Interactor usecases.ExecuteTransUsecase
}

// TransHandlerInput struct that represents the input
type TransHandlerInput struct {
	Command string                 `get:"command"`
	Params  map[string]interface{} `json:"params"`
	Format  string                 `json:"format"`
}

// TransRequestOutput struct that represents the output
type TransRequestOutput struct {
	Status   string            `json:"status"`
	Response map[string]string `json:"response"`
}

// TransRequestSliceOutput struct that represents the output slice
type TransRequestSliceOutput struct {
	Status   string              `json:"status"`
	Response []map[string]string `json:"response"`
}

// Input returns a fresh, empty instance of transHandlerInput
func (t *TransHandler) Input() HandlerInput {
	return &TransHandlerInput{}
}

// Execute executes the given trans request and returns the response
// of the execution.
// Expected response format:
//   { Status: string - "TRANS_OK" or error }
func (t *TransHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*TransHandlerInput)
	command := parseInput(in)
	transResp, err := t.Interactor.ExecuteCommand(command)
	// handle errors given by the interactor
	if err != nil {
		response = &goutils.Response{
			Code: http.StatusInternalServerError,
			Body: &goutils.GenericError{
				ErrorMessage: err.Error(),
			},
		}
		return response
	}

	// handle trans errors, database errors, or general reported errors by trans
	if transResp.Error() != nil ||
		transResp.Status() == usecases.TransError ||
		transResp.Status() == usecases.TransDatabaseError {
		mapResp, _ := transResp.Map()
		response = &goutils.Response{
			Code: http.StatusBadRequest,
			Body: TransRequestOutput{
				Status:   transResp.Status(),
				Response: mapResp,
			},
		}
		return response
	}

	switch in.Format {
	case "slice":
		sliceResp, _ := transResp.Slice()
		response = &goutils.Response{
			Code: http.StatusOK,
			Body: TransRequestSliceOutput{
				Status:   transResp.Status(),
				Response: sliceResp,
			},
		}
		return response
	default:
		mapRes, _ := transResp.Map()
		response = &goutils.Response{
			Code: http.StatusOK,
			Body: TransRequestOutput{
				Status:   transResp.Status(),
				Response: mapRes,
			},
		}
		return response
	}

}

func parseInput(input *TransHandlerInput) domain.TransCommand {
	command := domain.TransCommand{
		Command: input.Command,
	}

	params := make([]domain.TransParams, 0)

	for key, value := range input.Params {
		switch value.(type) {
		case []interface{}:
			for _, val := range value.([]interface{}) {
				switch val.(type) {
				case map[string]interface{}:
					for k, v := range val.(map[string]interface{}) {
						param := domain.TransParams{
							Key:   k,
							Value: v,
							Blob:  key == "blobs",
						}
						params = append(params, param)
					}
				case string:
					param := domain.TransParams{
						Key:   key,
						Value: val,
					}
					params = append(params, param)
				}
			}
		default:
			param := domain.TransParams{
				Key:   key,
				Value: value,
			}
			params = append(params, param)
		}
	}
	command.Params = params
	return command
}
