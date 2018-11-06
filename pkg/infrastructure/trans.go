package infrastructure

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.schibsted.io/Yapo/trans/pkg/interfaces/loggers"
	"github.schibsted.io/Yapo/trans/pkg/interfaces/repository/services"
)

// trans struct definition
type trans struct {
	conf            TransConf
	logger          loggers.Logger
	allowedCommands []string
}

// textProtocolTransFactory is a auxiliar struct to create trans on demand
type textProtocolTransFactory struct {
	conf            TransConf
	logger          loggers.Logger
	allowedCommands []string
}

// NewTextProtocolTransFactory initialize a services.TransFactory
func NewTextProtocolTransFactory(
	conf TransConf,
	logger loggers.Logger,
) services.TransFactory {
	return &textProtocolTransFactory{
		conf:            conf,
		logger:          logger,
		allowedCommands: strings.Split(conf.AllowedCommands, "|"),
	}
}

// MakeTransHandler initialize a services.TransHandler on demand
func (t *textProtocolTransFactory) MakeTransHandler() services.TransHandler {
	return &trans{
		conf:            t.conf,
		logger:          t.logger,
		allowedCommands: t.allowedCommands,
	}
}

// SendCommand use a socket conection to send commands to trans port
func (handler *trans) SendCommand(cmd string, params map[string]string) (map[string]string, error) {
	respMap := make(map[string]string)
	// check if the command is allowed; if not, return error
	valid := handler.isAllowedCommand(cmd)
	if !valid {
		err := fmt.Errorf(
			"Invalid Command. Valid commands: %s",
			handler.allowedCommands,
		)
		handler.logger.Debug(err.Error())
		return respMap, err
	}
	conn, err := handler.connect()
	defer conn.Close() //nolint: errcheck, megacheck

	if err != nil {
		handler.logger.Debug("Error connecting to trans: %s\n", err.Error())
		return respMap, err
	}
	// initiate the context so the request can timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(handler.conf.Timeout)*time.Second)
	defer cancel()

	respMap, err = handler.sendWithContext(ctx, conn, cmd, params)
	if err != nil {
		handler.logger.Debug("Error Sending command %s: %s\n", cmd, err)
	}

	return respMap, err
}

// isAllowedCommand checks if the given command can be sent to trans
func (handler *trans) isAllowedCommand(cmd string) bool {
	for _, allowedCommand := range handler.allowedCommands {
		if allowedCommand == cmd {
			return true
		}
	}
	return false
}

// connect returns a connection to the trans client.
// Retries to connect after retryAfter time if the connection times out
func (handler *trans) connect() (net.Conn, error) {
	// initiate the retrier that will handle retry reconnect if the connection dies
	r := retrier.New(
		[]time.Duration{
			time.Duration(handler.conf.RetryAfter) * time.Second},
		nil,
	)
	var conn net.Conn
	// set the function that starts the connection
	err := r.Run(func() error {
		var e error
		conn, e = net.DialTimeout(
			"tcp",
			fmt.Sprintf(
				"%s:%d",
				handler.conf.Host,
				handler.conf.Port,
			),
			time.Duration(handler.conf.Timeout)*time.Second,
		)
		return e
	})
	return conn, err
}

// sendWithContext sends the message to trans but is cancelable via a context.
// The context timeout specified how long the caller can wait
// for the trans to respond
func (handler *trans) sendWithContext(ctx context.Context, conn io.ReadWriteCloser, cmd string, args map[string]string) (map[string]string, error) {
	var resp map[string]string
	errChan := make(chan error, 1)

	// starts the go routine that sends the message and retrieves the response and error, if any.
	// it communicates any error through errChan
	go func() {
		errChan <- func() error {
			var err error
			resp, err = handler.send(conn, cmd, args)
			return err
		}()
	}()

	select {
	case <-ctx.Done():
		// closing the connection here interrupts the send function, in the gorouting above, if it
		// is waiting on reading from or writing to the connection.
		err := conn.Close()
		if err != nil {
			handler.logger.Debug("Error Closing connection to trans after ctx done: %s\n", err.Error())
		}
		// wait for the goroutine to return and ignore the error
		<-errChan
		// return the context error: the operation timed out.
		return nil, ctx.Err()
	case err := <-errChan:
		// in this case the send function returned before
		// the timeout of the context.
		return resp, err
	}
}

func (handler *trans) send(conn io.ReadWriter, cmd string, args map[string]string) (map[string]string, error) {
	// Check greeting.
	br := bufio.NewReaderSize(conn, handler.conf.BuffSize)
	line, err := br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(line, []byte("220 Welcome.\n")) {
		return nil, fmt.Errorf("trans: unexpected greeting: %q", line)
	}

	buf := make([]byte, 0)
	// Send command to Trans.
	buf = appendCmd(buf, cmd, args)
	if _, err = conn.Write(buf); err != nil {
		return nil, err
	}

	// Get response.
	buf = nil
	for {
		line, err = br.ReadSlice('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			if !bytes.HasSuffix(buf, []byte("end\n")) {
				return nil, fmt.Errorf("trans: response truncated: %q", line)
			}
			buf = buf[:len(buf)-4] // Remove "end".
			break
		}

		buf = append(buf, line...)
	}

	respMap, err := TransResponse(buf).Map()
	if err != nil {
		return respMap, fmt.Errorf("error parsing response: %s", err.Error())
	}
	return respMap, nil
}

// appendCmd Appends the command to the buffer. For the command format, see:
// https://scmcoord.com/wiki/Trans#Protocol
func appendCmd(buf []byte, cmd string, args map[string]string) []byte {
	buf = append(buf, "cmd:"...)
	buf = append(buf, cmd...)
	buf = append(buf, '\n')
	for key, value := range args {
		isBlob := isBlob(value)
		if isBlob {
			buf = append(buf, "blob:"...)
			buf = strconv.AppendInt(buf, int64(len(value)), 10)
			buf = append(buf, ':')
		}
		buf = append(buf, key...)
		if isBlob {
			buf = append(buf, '\n')
		} else {
			buf = append(buf, ':')
		}
		buf = append(buf, value...)
		buf = append(buf, '\n')
	}
	buf = append(buf, "commit:1"...)
	buf = append(buf, "\nend\n"...)
	return buf
}

//isBlob returns if the value is a blob (contains \n)
func isBlob(value string) bool {
	return strings.Contains(value, "\n")
}
