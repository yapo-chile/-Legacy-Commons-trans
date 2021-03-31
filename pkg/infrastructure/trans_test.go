package infrastructure

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.mpi-internal.com/Yapo/trans/pkg/domain"

	"github.com/stretchr/testify/assert"
	"github.mpi-internal.com/Yapo/trans/pkg/usecases"
)

func TestIsAllowedCommand(t *testing.T) {
	transHandler := trans{
		allowedCommands: []string{"transinfo", "get_account", "newad"},
	}

	assert.True(t, transHandler.isAllowedCommand("transinfo"))
	assert.False(t, transHandler.isAllowedCommand("loadad"))
	assert.True(t, transHandler.isAllowedCommand("get_account"))
	assert.False(t, transHandler.isAllowedCommand("Get_account"))
	assert.False(t, transHandler.isAllowedCommand("newAd"))
	assert.False(t, transHandler.isAllowedCommand("newad:"))
	assert.True(t, transHandler.isAllowedCommand("newad"))
}

func TestSendCommandInvalidCommand(t *testing.T) {
	//initiate the conf
	host := "" //shouldn't try to connect with the server
	port := 0
	conf := TransConf{
		Host:            host,
		Port:            port,
		Timeout:         15,
		RetryAfter:      5,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	logger.On("Error")
	expectedResponse := make([]map[string]string, 1)
	expectedResponse[0] = map[string]string{"error": "Invalid Command. Valid commands: [test]"}
	cmd := "transinfo"
	params := []domain.TransParams{
		{
			Key:   "param1",
			Value: "ok",
		},
	}

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()

	resp, err := transHandler.SendCommand(cmd, params)
	assert.Error(t, err)
	assert.Equal(t, expectedResponse, resp)
	logger.AssertExpectations(t)
}

func TestSendCommandTimeout(t *testing.T) {
	command := "cmd:test\nparam1:ok\ncommit:1\nend\n"
	response := fmt.Sprintf("status:%s\n", usecases.TransOK)
	//define the function that will receive the message
	handlerFunc := func(input []byte) []byte {
		time.Sleep(2 * time.Second)
		// in case the request reaches after the sleep
		assert.Equal(t, command, string(input))
		return []byte(response)
	}
	//initiate the mock server
	server := NewMockTransServer()
	defer server.Close()
	server.SetHandler(handlerFunc)

	//initiate the conf
	addr := strings.Split(server.Address, ":")
	host := addr[0]
	port, _ := strconv.Atoi(addr[1])
	conf := TransConf{
		Host:            host,
		Port:            port,
		Timeout:         1,
		RetryAfter:      5,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	logger.On("Error")
	var expectedResponse []map[string]string
	cmd := "test"
	params := []domain.TransParams{
		domain.TransParams{
			Key:   "param1",
			Value: "ok",
		},
	}

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()
	resp, err := transHandler.SendCommand(cmd, params)
	assert.Error(t, err)
	assert.Equal(t, expectedResponse, resp)
	logger.AssertExpectations(t)
}

func TestSendCommandBusyServer(t *testing.T) {
	//initiate the mock server
	server := NewMockTransServer()
	defer server.Close()
	server.SetBusy(true)

	//initiate the conf
	addr := strings.Split(server.Address, ":")
	host := addr[0]
	port, _ := strconv.Atoi(addr[1])
	conf := TransConf{
		Host:            host,
		Port:            port,
		Timeout:         15,
		RetryAfter:      5,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	logger.On("Error")
	var expectedResponse []map[string]string
	cmd := "test"
	params := []domain.TransParams{
		domain.TransParams{
			Key:   "param1",
			Value: "ok",
		},
	}

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()

	resp, err := transHandler.SendCommand(cmd, params)
	assert.Error(t, err)
	assert.Equal(t, expectedResponse, resp)
	logger.AssertExpectations(t)
}

func TestSendCommandOK(t *testing.T) {

	command := "cmd:test\nparam1:ok\xc1\ncommit:1\nend\n"
	response := "status:TRANS_OK\n"

	//define the function that will receive the message
	handlerFunc := func(input []byte) []byte {
		assert.Equal(t, command, string(input))
		return []byte(response)
	}
	//initiate the server
	server := NewMockTransServer()
	defer server.Close()
	server.SetHandler(handlerFunc)

	//initiate the conf
	addr := strings.Split(server.Address, ":")
	host := addr[0]
	port, _ := strconv.Atoi(addr[1])
	conf := TransConf{
		Host:            host,
		Port:            port,
		Timeout:         15,
		RetryAfter:      5,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	expectedResponse := make([]map[string]string, 1)
	expectedResponse[0] = map[string]string{"status": usecases.TransOK}
	cmd := "test"
	params := []domain.TransParams{
		domain.TransParams{
			Key:   "param1",
			Value: "okÁ",
		},
	}

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()

	resp, err := transHandler.SendCommand(cmd, params)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
	logger.AssertExpectations(t)
}

func TestSendCommandBlobOK(t *testing.T) {

	command := "cmd:test\nblob:5:body\nedgar\ncommit:1\nend\n"
	response := "status:TRANS_OK\n"

	//define the function that will receive the message
	handlerFunc := func(input []byte) []byte {
		assert.Equal(t, command, string(input))
		return []byte(response)
	}
	//initiate the server
	server := NewMockTransServer()
	defer server.Close()
	server.SetHandler(handlerFunc)

	//initiate the conf
	addr := strings.Split(server.Address, ":")
	host := addr[0]
	port, _ := strconv.Atoi(addr[1])
	conf := TransConf{
		Host:            host,
		Port:            port,
		Timeout:         15,
		RetryAfter:      5,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	expectedResponse := make([]map[string]string, 1)
	expectedResponse[0] = map[string]string{"status": usecases.TransOK}
	cmd := "test"
	params := []domain.TransParams{
		domain.TransParams{
			Key:   "body",
			Value: "ZWRnYXI=",
			Blob:  true,
		},
	}

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()

	resp, err := transHandler.SendCommand(cmd, params)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
	logger.AssertExpectations(t)
}
func TestISO8859Input(t *testing.T) {
	//define the function that will receive the message
	handlerFunc := func(input []byte) []byte {
		var response []byte
		return []byte(response)
	}
	//initiate the server
	server := NewMockTransServer()
	defer server.Close()
	server.SetHandler(handlerFunc)

	//initiate the conf
	addr := strings.Split(server.Address, ":")
	host := addr[0]
	port, _ := strconv.Atoi(addr[1])
	conf := TransConf{
		Host:            host,
		Port:            port,
		Timeout:         15,
		RetryAfter:      5,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	cmd := "test"
	params := []domain.TransParams{
		{
			Key:   "param1",
			Value: "ok\xc1",
		},
	}

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()

	resp, err := transHandler.SendCommand(cmd, params)
	assert.NoError(t, err)
	assert.Equal(t, []map[string]string{}, resp)
	logger.AssertExpectations(t)
}
