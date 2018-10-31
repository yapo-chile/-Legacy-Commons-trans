package infrastructure

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendCommandOK(t *testing.T) {
	command := "cmd:test\nparam1:ok\ncommit:1\nend\n"
	response := "status:TRANS_OK\nend\n"
	//initiate a tcp listener for the test
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		defer ln.Close()
		conn, _ := ln.Accept()
		io.WriteString(conn, "220 Welcome.\n")
		br := bufio.NewReader(conn)
		buf := make([]byte, 0, 100)
		for {
			line, _ := br.ReadSlice('\n')
			buf = append(buf, line...)
			if bytes.Equal(line, []byte("end\n")) {
				break
			}
		}
		receivedCommand := string(buf)
		assert.Equal(t, command, receivedCommand)
		conn.Write([]byte(response))
		conn.Close()
	}()

	addr := strings.Split(ln.Addr().String(), ":")
	host := addr[0]
	port, _ := strconv.Atoi(addr[1])
	conf := TransConf{
		Host:       host,
		Port:       port,
		Timeout:    15,
		RetryAfter: 5,
		BuffSize:   4096,
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
