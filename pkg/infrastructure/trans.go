package infrastructure

import (
	"fmt"
	"io/ioutil"
	"net"
	"sort"
	"strings"

	"github.schibsted.io/Yapo/trans/pkg/interfaces/loggers"
	"github.schibsted.io/Yapo/trans/pkg/interfaces/repository/services"
)

// Trans struct definition
type Trans struct {
	Conf   TransConf
	Logger loggers.Logger
}

// TextProtocolTransFactory is a auxiliar struct to create trans on demand
type TextProtocolTransFactory struct {
	Conf   TransConf
	Logger loggers.Logger
}

// NewTextProtocolTransFactory initialize a TextProtocolTransFactory
func NewTextProtocolTransFactory(
	conf TransConf,
	logger loggers.Logger,
) *TextProtocolTransFactory {
	return &TextProtocolTransFactory{
		Conf:   conf,
		Logger: logger,
	}
}

// MakeTransHandler initialize a TransHandler on demand
func (t *TextProtocolTransFactory) MakeTransHandler() services.TransHandler {
	return &Trans{
		Conf:   t.Conf,
		Logger: t.Logger,
	}
}

// SendCommand use a socket conection to send commands to trans port
func (handler *Trans) SendCommand(cmd string, params map[string]string) (map[string]string, error) {
	respMap := make(map[string]string)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", handler.Conf.Host, handler.Conf.Port))
	if err != nil {
		handler.Logger.Debug("Error Trans Conection: %s\n", err)
		return respMap, err
	}

	commandString := fmt.Sprintf("cmd:%s\n%s\nend\n", cmd, handler.parseFormatTrans(params, true))
	_, err = fmt.Fprint(conn, commandString)
	if err != nil {
		handler.Logger.Debug("Error Sending command %s: %s\n", cmd, err)
		return respMap, err
	}
	// get the response and check if there are errors
	resp, err := ioutil.ReadAll(conn)
	stringResp := fmt.Sprintf("%s", resp)
	// cast the response into a map
	splitResp := strings.Split(stringResp, "\n")
	for _, pair := range splitResp {
		splitPair := strings.Split(pair, ":")
		if len(splitPair) > 1 {
			respMap[splitPair[0]] = splitPair[1]
		}
	}
	if status, ok := respMap["status"]; ok {
		if status == "TRANS_ERROR" {
			handler.Logger.Debug("\n---- trans response start ----\n%s---- trans response end ----\n", stringResp)
			errorMessage := "Error in TRANS response"
			if errString, ok := respMap["error"]; ok {
				errorMessage = errString
			}
			return respMap, fmt.Errorf("%s", errorMessage)
		}
	}
	return respMap, err
}

// parseFormatTrans parse a mapping string to trans string format
func (handler *Trans) parseFormatTrans(m map[string]string, commit bool) (s string) {
	// sort by string for pact-test
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if strings.ContainsAny(m[k], "\n") {
			handler.Logger.Error("Trans does not accept character newline in value: %s", m[k])
			return ""
		}
		s += fmt.Sprintf("%s:%s\n", k, m[k])
	}
	if commit {
		s += "commit:1\n"
	}
	return
}
