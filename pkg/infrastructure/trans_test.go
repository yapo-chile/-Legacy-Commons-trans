package infrastructure

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestIsBlob(t *testing.T) {
	assert.True(t, isBlob("key\nvalue"))
	assert.False(t, isBlob("key\\nvalue"))
	assert.True(t, isBlob("value\n"))
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
		BuffSize:        4096,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	logger.On("Debug")
	expectedResponse := make(map[string]string)
	cmd := "transinfo"
	params := make(map[string]string)
	params["param1"] = "ok"

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()

	resp, err := transHandler.SendCommand(cmd, params)
	assert.Error(t, err)
	assert.Equal(t, expectedResponse, resp)
}

func TestSendCommandTimeout(t *testing.T) {
	command := "cmd:test\nparam1:ok\ncommit:1\nend\n"
	response := "status:TRANS_OK\n"
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
		BuffSize:        4096,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	logger.On("Debug")
	var expectedResponse map[string]string
	cmd := "test"
	params := make(map[string]string)
	params["param1"] = "ok"

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()
	resp, err := transHandler.SendCommand(cmd, params)
	assert.Error(t, err)
	assert.Equal(t, expectedResponse, resp)
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
		BuffSize:        4096,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	logger.On("Debug")
	var expectedResponse map[string]string
	cmd := "test"
	params := make(map[string]string)
	params["param1"] = "ok"

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()

	resp, err := transHandler.SendCommand(cmd, params)
	assert.Error(t, err)
	assert.Equal(t, expectedResponse, resp)
}

func TestSendCommandOK(t *testing.T) {
	command := "cmd:test\nparam1:ok\ncommit:1\nend\n"
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
		BuffSize:        4096,
		AllowedCommands: "test",
	}
	logger := MockLoggerInfrastructure{}
	logger.On("Debug")
	expectedResponse := make(map[string]string)
	expectedResponse["status"] = "TRANS_OK"
	cmd := "test"
	params := make(map[string]string)
	params["param1"] = "ok"

	transFactory := NewTextProtocolTransFactory(conf, &logger)
	transHandler := transFactory.MakeTransHandler()

	resp, err := transHandler.SendCommand(cmd, params)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
}
