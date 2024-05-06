package webserver

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net"
	"os"
)

type Message struct {
	Endpoint string
	Body     string
}

func (m *Message) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

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
	socketDir := homeDir + "/liScreMon/socket/liScreMon.sock"

	if unixAddr, err = net.ResolveUnixAddr(unix, socketDir); err != nil {
		return nil, err
	}

	conn, err := net.DialUnix(unix, nil, unixAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (a *App) CheckDaemonService() error {
	msg := Message{
		Endpoint: "socketAlive",
		Body:     "I wish this project prospered.",
	}
	bytes, err := msg.encode()
	if err != nil {
		return err
	}
	_, err = a.daemonConn.Write(bytes)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	n, err := a.daemonConn.Read(buf)
	if err != nil {
		return err
	}

	response := string(buf[:n])
	fmt.Printf("string message from daemon----> \"%s\"\n", response)
	return nil
}
