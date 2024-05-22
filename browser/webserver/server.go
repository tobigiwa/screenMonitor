package webserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	helperFuncs "pkg/helper"
	"pkg/types"
)

type App struct {
	logger     *slog.Logger
	daemonConn net.Conn
}

func NewApp(logger *slog.Logger) (*App, error) {
	daemonConn, err := listenToDaemonService()
	if err != nil {
		return nil, err
	}

	return &App{
		logger:     logger,
		daemonConn: daemonConn,
	}, nil
}

func listenToDaemonService() (net.Conn, error) {
	var (
		unix     = "unix"
		homeDir  string
		err      error
		unixAddr *net.UnixAddr
	)
	if homeDir, err = os.UserHomeDir(); err != nil {
		return nil, err
	}
	socketDir := homeDir + "/liScreMon/socket/daemon.sock"

	if unixAddr, err = net.ResolveUnixAddr(unix, socketDir); err != nil {
		return nil, err
	}

	conn, err := net.DialUnix(unix, nil, unixAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (a *App) CheckDaemonService() (types.Message, error) {
	msg := types.Message{
		Endpoint:          "startConnection",
		StringDataRequest: "I wish this project prospered.",
	}
	return a.writeAndReadWithDaemonService(msg)
}

func (a *App) writeAndReadWithDaemonService(msg types.Message) (types.Message, error) {
	bytesData, err := helperFuncs.Encode(msg) // encode message in byte
	if err != nil {
		return types.NoMessage, fmt.Errorf("encode %w", err)
	}
	if _, err = a.daemonConn.Write(bytesData); err != nil { // write to socket
		return types.NoMessage, fmt.Errorf("write %w", err)
	}

	var dataBuf bytes.Buffer
	tempBuf := make([]byte, 100_000) //100kb
	n := 0

	for {
		if n, err = a.daemonConn.Read(tempBuf); err != nil { // read response from socket
			if errors.Is(err, io.EOF) {
				break
			}
			return types.NoMessage, fmt.Errorf("read error from socket %w", err)
		}

		if n > 0 {
			dataBuf.Write(tempBuf[:n])
		}

		if json.Valid(dataBuf.Bytes()) { // Implement this function based on your protocol
			break
		}
	}

	return helperFuncs.Decode[types.Message](dataBuf.Bytes())
}
