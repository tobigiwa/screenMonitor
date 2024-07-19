package backend

import (
	"context"
	"log/slog"

	utils "utils"
)

func (a *App) CloseDaemonConnection() error {

	msg := utils.Message{
		Endpoint: "closeConnection",
	}

	bytes, err := utils.EncodeJSON(msg)
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
