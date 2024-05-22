package webserver

import (
	"context"
	"log/slog"
	helperFuncs "pkg/helper"
	"pkg/types"
)

func (a *App) CloseDaemonConnection() error {

	msg := types.Message{
		Endpoint: "closeConnection",
	}

	bytes, err := helperFuncs.Encode(msg)
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
