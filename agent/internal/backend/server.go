package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"path/filepath"
	"runtime"

	utils "utils"
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
		unix      = "unix"
		socketDir string
		err       error
		unixAddr  *net.UnixAddr
	)

	if runtime.GOOS == "linux" {
		socketDir = filepath.Join(utils.APP_CONFIG_DIR, "socket", "daemon.sock")
	}
	if runtime.GOOS == "windows" {
		notImplemented()
	}
	if runtime.GOOS == "darwin" {
		notImplemented()
	}

	if unixAddr, err = net.ResolveUnixAddr(unix, socketDir); err != nil {
		return nil, err
	}

	conn, err := net.DialUnix(unix, nil, unixAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (a *App) CheckDaemonService() (utils.Message, error) {
	msg := utils.Message{
		Endpoint:    "startConnection",
		StatusCheck: "I wish this project prospered.",
	}
	return a.commWithDaemonService(msg)
}

func (a *App) commWithDaemonService(msg utils.Message) (utils.Message, error) {
	bytesData, err := utils.EncodeJSON(msg) // encode message in byte
	if err != nil {
		return utils.NoMessage, fmt.Errorf("encode %w", err)
	}

	if _, err = a.daemonConn.Write(bytesData); err != nil { // write to socket
		return utils.NoMessage, fmt.Errorf("write %w", err)
	}

	var dataBuf bytes.Buffer
	tempBuf := make([]byte, 100_000) //100kb
	n := 0

	for {
		if n, err = a.daemonConn.Read(tempBuf); err != nil { // read response from socket
			if errors.Is(err, io.EOF) {
				break
			}
			return utils.NoMessage, fmt.Errorf("read error from socket %w", err)
		}

		if n > 0 {
			dataBuf.Write(tempBuf[:n])
		}

		if json.Valid(dataBuf.Bytes()) {
			break
		}
	}

	if msg, err = utils.DecodeJSON[utils.Message](dataBuf.Bytes()); err != nil {
		return utils.NoMessage, err
	}

	if msg.IsError {
		return utils.NoMessage, fmt.Errorf(msg.Error)
	}
	return msg, nil
}

func notImplemented() {}
