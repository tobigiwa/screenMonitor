package webserver

import (
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
	bytes, err := helperFuncs.Encode(msg)
	if err != nil {
		return types.NoMessage, err
	}
	if _, err = a.daemonConn.Write(bytes); err != nil {
		return types.NoMessage, err
	}
	buf := make([]byte, 10240)
	if _, err := a.daemonConn.Read(buf); err != nil {
		return types.NoMessage, err
	}
	return helperFuncs.Decode[types.Message](buf)
}
