package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
)

type Command byte
type ResponseStatus byte

const (
	CommandCopy Command = iota
	CommandPaste
)

const (
	ResponseStatusOK ResponseStatus = iota

	// ResponseStatusErrInternalErr means a generic internal error occurred on the server.
	ResponseStatusErrInternalErr

	// ResponseStatusErrSizeMismatch means client did not send as much data as expected.
	ResponseStatusErrSizeMismatch

	// ResponseStatusErrClipboardWriteErr means server was unable to write to clipboard for some internal reason.
	ResponseStatusErrClipboardWriteErr

	// ResponseStatusErrClipboardReadErr means server was unable to read from clipboard for some internal reason.
	ResponseStatusErrClipboardReadErr

	// ResponseStatusErrInvalidCommand means client sent an invalid command.
	ResponseStatusErrInvalidCommand

	// ResponseStatusErrBadRequest means client sent a request that was somehow invalid.
	ResponseStatusErrBadRequest
)

var ResponseStatusMessages = map[ResponseStatus]string{
	ResponseStatusOK:                   "OK",
	ResponseStatusErrInternalErr:       "Internal server error",
	ResponseStatusErrSizeMismatch:      "Size mismatch",
	ResponseStatusErrClipboardWriteErr: "Failed to write to clipboard",
	ResponseStatusErrClipboardReadErr:  "Failed to read from clipboard",
	ResponseStatusErrInvalidCommand:    "Invalid command",
	ResponseStatusErrBadRequest:        "Bad request",
}

type ClientRequest struct {
	Command Command
	Data    string
}

func (c *ClientRequest) Bytes() []byte {
	buf := make([]byte, 9)
	buf[0] = byte(c.Command)
	binary.BigEndian.PutUint64(buf[1:], uint64(len(c.Data)))
	return append(buf, []byte(c.Data)...)
}

func (c *ClientRequest) Read(conn net.Conn) (ResponseStatus, error) {
	buf := make([]byte, 9)
	n, err := conn.Read(buf)
	if err != nil {
		return ResponseStatusErrInternalErr, fmt.Errorf("failed to read from connection: %w", err)
	}

	if n != 9 {
		return ResponseStatusErrBadRequest, fmt.Errorf("invalid message, expected 9 bytes but got %d", n)
	}

	c.Command = Command(buf[0])
	size := binary.BigEndian.Uint64(buf[1:])

	if c.Command != CommandCopy && c.Command != CommandPaste {
		return ResponseStatusErrInvalidCommand, fmt.Errorf("invalid command: %d", c.Command)
	}

	if c.Command == CommandPaste && size > 0 {
		return ResponseStatusErrBadRequest, fmt.Errorf("client sent %d bytes for paste command, expected 0", size)
	}

	if size > 0 {
		data := make([]byte, size)
		n, err := conn.Read(data)
		if err != nil {
			return ResponseStatusErrInternalErr, fmt.Errorf("failed to read from connection: %w", err)
		}

		if uint64(n) != size {
			return ResponseStatusErrSizeMismatch, fmt.Errorf("client sent %d bytes, expected %d", n, size)
		}

		c.Data = string(data)
	}

	return ResponseStatusOK, nil
}

type ServerResponse struct {
	Status ResponseStatus
	Data   string
}

func (r *ServerResponse) Bytes() []byte {
	buf := make([]byte, 9)
	buf[0] = byte(r.Status)
	binary.BigEndian.PutUint64(buf[1:], uint64(len(r.Data)))
	return append(buf, r.Data...)
}

func (r *ServerResponse) Read(conn net.Conn) error {
	buf := make([]byte, 9)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read from socket: %w", err)
	}

	if n != 9 {
		return fmt.Errorf("invalid message, expected 9 bytes but got %d", n)
	}

	r.Status = ResponseStatus(buf[0])
	size := binary.BigEndian.Uint64(buf[1:])

	dataBuf := make([]byte, size)
	n, err = conn.Read(dataBuf)
	if err != nil {
		return fmt.Errorf("failed to read from socket: %w", err)
	}

	if uint64(n) != size {
		return fmt.Errorf("server sent %d bytes, expected %d", n, size)
	}

	r.Data = string(dataBuf)
	return nil
}

func Serve(socketFile string) error {
	socket, err := net.Listen("unix", socketFile)
	if err != nil {
		return fmt.Errorf("unable to listen on unix socket '%s': %w", socketFile, err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		_ = socket.Close()
		_ = os.Remove(socketFile)
		os.Exit(0)
	}()

	slog.Info(fmt.Sprintf("Server listening on '%s'", socketFile))

	for {
		conn, err := socket.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				slog.Error(fmt.Sprintf("Failed to accept connection: %s", err))
			}
			continue
		}

		if err := conn.SetDeadline(time.Now().Add(time.Second * 10)); err != nil {
			slog.Error(fmt.Sprintf("Failed to set deadline on connection: %s", err))
			continue
		}

		slog.Debug("Accepted connection")
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	request := &ClientRequest{}
	response := &ServerResponse{}

	responseStatus, err := request.Read(conn)
	if err != nil {
		slog.Error(err.Error())
		response.Status = responseStatus
		response.Data = ""
		_, err := conn.Write(response.Bytes())
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to write to connection: %s", err))
		}

		return
	}

	switch request.Command {
	case CommandCopy:
		if err := clipboard.WriteAll(request.Data); err != nil {
			errMsg := fmt.Sprintf("Failed to write to clipboard: %s", err)
			slog.Error(errMsg)
			response.Status = ResponseStatusErrClipboardWriteErr
			response.Data = errMsg
			_, err := conn.Write(response.Bytes())
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to write to connection: %s", err))
			}

			return
		}

		response.Status = ResponseStatusOK
		response.Data = ""
		_, err = conn.Write(response.Bytes())
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to write to connection: %s", err))
		}

		slog.Info(fmt.Sprintf("Copied %d bytes to clipboard", len(request.Data)))
	case CommandPaste:
		text, err := clipboard.ReadAll()
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read from clipboard: %s", err))
			response.Status = ResponseStatusErrClipboardReadErr
			response.Data = ""
			_, err := conn.Write(response.Bytes())
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to write to connection: %s", err))
			}

			return
		}

		response.Status = ResponseStatusOK
		response.Data = text
		_, err = conn.Write(response.Bytes())
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to write to connection: %s", err))
		}

		slog.Info(fmt.Sprintf("Pasted %d bytes from clipboard", len(text)))
	}
}
