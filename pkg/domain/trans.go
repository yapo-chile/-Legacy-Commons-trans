package domain

// TransParams is a struct with Trans format params
type TransParams struct {
	Key   string
	Value interface{}
	Blob  bool
}

// TransCommand represents a trans command with params to be executed on a trans server
type TransCommand struct {
	// the command to be executed
	Command string
	// Params the params of the command
	Params []TransParams
}

// TransResponse represents the response given to the execution of a TransCommand
type TransResponse interface {
	// Status the status of the response (normally TRANS_OK or TRANS_ERROR)
	Status() string
	SetStatus(string)
	// Map returns response as map format
	Map() (map[string]string, error)
	// Slice returns response as slice format
	Slice() ([]map[string]string, error)
	Error() error
	SetError(error)
}

// TransRepository defines a storage for the trans commands
type TransRepository interface {
	// Execute executes the command on a trans server
	Execute(command TransCommand) (TransResponse, error)
}
