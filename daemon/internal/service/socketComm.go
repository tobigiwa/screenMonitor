package service

import (
	db "LiScreMon/daemon/internal/database"
	"LiScreMon/daemon/internal/tasks"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	helperFuncs "pkg/helper"
	"pkg/types"
	"syscall"
)

func NewService(db *db.BadgerDBStore) (*Service, error) {

	var (
		service Service
		err     error
	)

	service.db = db

	if service.taskManager, err = tasks.StartTaskManger(db); err != nil {
		return nil, fmt.Errorf("error starting task manager")
	}

	return &service, nil
}

func (s *Service) StartService(socketDir string, db *db.BadgerDBStore) error {

	SocketConn, err := domainSocket(socketDir)
	if err != nil {
		return err
	}
	s.handleConnection(SocketConn) // blocking

	SocketConn.Close()
	return nil
}

func domainSocket(socketDir string) (*net.UnixListener, error) {

	var (
		conn     *net.UnixListener
		err      error
		unixAddr *net.UnixAddr
	)

	if err = os.MkdirAll(socketDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating socket dir:%w", err)
	}

	socketFilePath := filepath.Join(socketDir, "daemon.sock")

	syscall.Unlink(socketFilePath)

	if unixAddr, err = net.ResolveUnixAddr("unix", socketFilePath); err != nil {
		return nil, err
	}

	if conn, err = net.ListenUnix("unix", unixAddr); err != nil {
		return nil, fmt.Errorf("error creating domain socket:%w", err)
	}

	conn.SetUnlinkOnClose(true)

	return conn, nil
}

func (s *Service) handleConnection(listener *net.UnixListener) {
	for {
		c, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Println("connection closed, daemonService says 🖐")
				return
			} else {
				fmt.Println("error accepting connection...buh we muuve:", err)
				continue
			}
		}

		fmt.Println("Connection accepted")
		go s.treatMessage(c)
	}
}

func (s *Service) treatMessage(c net.Conn) {
	for {
		var (
			msg types.Message
			err error
			n   int
		)

		buf := make([]byte, 10_000) //10kb
		if n, err = c.Read(buf); err != nil {
			fmt.Println("error reading message:", err)
			if errors.Is(err, io.EOF) {
				fmt.Println("client connection closed")
				c.Close()
				return
			}
			continue
		}

		if msg, err = helperFuncs.DecodeJSON[types.Message](buf[:n]); err != nil {
			fmt.Println("error decoding socket message", err)
			c.Close()
			return
		}

		switch msg.Endpoint {

		case "startConnection":
			msg = types.Message{StatusCheck: "OK"}

		case "closeConnection":
			closeConnection(c)
			return

		case "weekStat":
			msg.WeekStatResponse, err = s.getWeekStat(msg.WeekStatRequest)

		case "appStat":
			msg.AppStatResponse, err = s.getAppStat(msg.AppStatRequest)

		case "dayStat":
			msg.DayStatResponse, err = s.getDayStat(msg.DayStatRequest)

		case "setCategory":
			msg.SetCategoryResponse, err = s.setAppCategory(msg.SetCategoryRequest)

		case "tasks":
			msg.ReminderAndLimitResponse, err = s.tasks()

		case "reminders":
			msg.ReminderAndLimitResponse, err = s.allReminderTask()

		case "appLimits":
			msg.ReminderAndLimitResponse, err = s.allDailyAppLimitTask()

		case "newReminder":
			msg.ReminderAndLimitResponse, err = s.addNewReminder(msg.ReminderAndLimitRequest)

		case "newAppLimit":
			msg.ReminderAndLimitResponse, err = s.addNewLimitApp(msg.ReminderAndLimitRequest)

		case "removeTask":
			msg.ReminderAndLimitResponse, err = s.removeTask(msg.ReminderAndLimitRequest)
		}

		if err != nil {
			msg.IsError, msg.Error = true, err.Error()
		}

		bytes, err := helperFuncs.EncodeJSON(msg)
		if err != nil {
			msg.Error = fmt.Sprintf("error encoding response in serviceInstance: %s: %s", msg.Error, err.Error())
		}

		_, err = c.Write(bytes)
		if err != nil {
			fmt.Println("error writing response:", err)
			continue
		}
	}
}

func closeConnection(c net.Conn) {
	fmt.Println("we got a close connection message")
	err := c.Close()
	if err != nil {
		fmt.Println("ERROR CLOSING CONNECTION")
	}
}
