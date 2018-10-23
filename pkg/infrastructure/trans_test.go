package infrastructure

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendCommand(t *testing.T) {
	command := "cmd:test\nparam1:ok\ncommit:1\n\nend\n"
	response := "status:TRANS_OK\n"
	//initiate a tcp listener for the test
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		defer ln.Close()
		conn, _ := ln.Accept()
		receivedCommand := strings.Builder{}
		data := ""
		reader := bufio.NewReader(conn)
		for strings.Compare(data, "end\n") != 0 {
			data, _ = reader.ReadString('\n')
			receivedCommand.WriteString(data)
		}

		assert.Equal(t, command, receivedCommand.String())
		conn.Write([]byte(response))
		conn.Close()
	}()

	addr := strings.Split(ln.Addr().String(), ":")
	host := addr[0]
	port, _ := strconv.Atoi(addr[1])
	conf := TransConf{
		Host: host,
		Port: port,
	}
	logger := MockLoggerInfrastructure{}

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
