package webserver

import (
	"context"
	"log/slog"
)

func (a *App) CloseDaemonConnection() error {

	msg := Message{
		Endpoint: "closeConnection",
	}

	bytes, err := msg.encode()
	if err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return err
	}
	if _, err = a.daemonConn.Write(bytes); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return err
	}
	return a.daemonConn.Close()
}
