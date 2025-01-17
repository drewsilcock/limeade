package client

import (
	"fmt"
	"github.com/drewsilcock/limeade/server"
	"log/slog"
	"net"
	"time"
)

type Client struct {
	socketFile string
}

func New(socketFile string) *Client {
	return &Client{
		socketFile: socketFile,
	}
}

func (c *Client) Paste() (string, error) {
	slog.Debug("Requesting paste from server")

	conn, err := net.Dial("unix", c.socketFile)
	if err != nil {
		err = fmt.Errorf("unable to connect to socket: %w", err)
		slog.Error(err.Error())
		return "", err
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	if err := conn.SetDeadline(time.Now().Add(time.Second * 10)); err != nil {
		err = fmt.Errorf("failed to set deadline on connection: %w", err)
		slog.Error(err.Error())
		return "", err
	}

	request := server.ClientRequest{}
	request.Command = server.CommandPaste
	request.Data = ""

	if _, err := conn.Write(request.Bytes()); err != nil {
		err = fmt.Errorf("failed to write to connection: %w", err)
		slog.Error(err.Error())
		return "", err
	}

	response := &server.ServerResponse{}
	if err := response.Read(conn); err != nil {
		return "", err
	}

	return response.Data, nil
}

func (c *Client) Copy(text string) error {
	slog.Debug(fmt.Sprintf("Sending copy to clipboard: %s", text))

	conn, err := net.Dial("unix", c.socketFile)
	if err != nil {
		err = fmt.Errorf("unable to connect to socket: %w", err)
		slog.Error(err.Error())
		return err
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	if err := conn.SetDeadline(time.Now().Add(time.Second * 10)); err != nil {
		err = fmt.Errorf("failed to set deadline on connection: %w", err)
		slog.Error(err.Error())
		return err
	}

	request := server.ClientRequest{}
	request.Command = server.CommandCopy
	request.Data = text

	if _, err := conn.Write(request.Bytes()); err != nil {
		err = fmt.Errorf("failed to write to connection: %w", err)
		slog.Error(err.Error())
		return err
	}

	response := &server.ServerResponse{}
	if err := response.Read(conn); err != nil {
		err = fmt.Errorf("failed to read response: %w", err)
		slog.Error(err.Error())
		return err
	}

	if response.Status != server.ResponseStatusOK {
		err := fmt.Errorf(
			"server returned error status %s: %s",
			server.ResponseStatusMessages[response.Status],
			response.Data,
		)
		slog.Error(err.Error())
		return err
	}

	return nil
}
