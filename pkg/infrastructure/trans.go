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

const (
	// WelcomeMessage ...
	WelcomeMessage = "220 Welcome.\n"
	// EndMessage is command that define end of the message.
	EndMessage = "end\n"
)

// trans struct definition
type trans struct {
	Conf   TransConf
	Logger loggers.Logger
}

// textProtocolTransFactory is a auxiliar struct to create trans on demand
type textProtocolTransFactory struct {
	Conf   TransConf
	Logger loggers.Logger
}

// NewTextProtocolTransFactory initialize a TextProtocolTransFactory
func NewTextProtocolTransFactory(
	conf TransConf,
	logger loggers.Logger,
) services.TransFactory {
	return &textProtocolTransFactory{
		Conf:   conf,
		Logger: logger,
	}
}

// MakeTransHandler initialize a TransHandler on demand
func (t *textProtocolTransFactory) MakeTransHandler() services.TransHandler {
	return &trans{
		Conf:   t.Conf,
		Logger: t.Logger,
	}
}

// SendCommand use a socket conection to send commands to trans port
func (handler *trans) SendCommand(cmd string, params map[string]string) (map[string]string, error) {
	respMap := make(map[string]string)
	conn, err := handler.connect()
	if err != nil {
		handler.Logger.Debug("Error connecting to trans: %s\n", err.Error())
		return respMap, err
	}
	// defer the close of the connection
	defer func() {
		err = conn.Close()
		if err != nil {
			handler.Logger.Debug("Error Closing connection to trans: %s\n", err.Error())
		}
	}()

	// initiate the context so the request can timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(handler.Conf.Timeout)*time.Second)
	defer cancel()

	respMap, err = handler.sendWithContext(ctx, conn, cmd, params)
	if err != nil {
		handler.Logger.Debug("Error Sending command %s: %s\n", cmd, err)
	}

	return respMap, err
}

// connect returns a connection to the trans client.
// Retries to connect after retryAfter time if the connection times out
func (handler *trans) connect() (net.Conn, error) {

	r := retrier.New([]time.Duration{time.Duration(handler.Conf.RetryAfter) * time.Second}, nil)

	var conn net.Conn

	err := r.Run(func() error {
		var e error
		conn, e = net.DialTimeout(
			"tcp",
			fmt.Sprintf(
				"%s:%d",
				handler.Conf.Host,
				handler.Conf.Port,
			),
			time.Duration(handler.Conf.Timeout)*time.Second,
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
			handler.Logger.Debug("Error Closing connection to trans after ctx done: %s\n", err.Error())
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
	br := bufio.NewReaderSize(conn, handler.Conf.BuffSize)
	line, err := br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	line = line[:len(line)-1] // Strip newline.
	if !bytes.Equal(line, []byte("220 Welcome.")) {
		return nil, fmt.Errorf("trans: unexpected greeting: %q", line)
	}

	buf := make([]byte, 0, 100)

	// Send command to Trans.
	buf = appendCmd(buf, cmd, args)
	if _, err = conn.Write(buf); err != nil {
		return nil, err
	}

	// Get response.
	buf = buf[:0] // Reset buffer to reuse it to write the response.
	for {
		line, err = br.ReadSlice('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			if !bytes.HasSuffix(buf, []byte(EndMessage)) {
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
