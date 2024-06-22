package service

import (
	db "LiScreMon/cli/daemon/internal/database"
	"LiScreMon/cli/daemon/internal/jobs"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	helperFuncs "pkg/helper"
	"pkg/types"
	"syscall"
)

var (
	ServiceInstance Service
	SocketConn      *net.UnixListener
)

func StartService(socketDir string, db *db.BadgerDBStore) {

	ServiceInstance.db = db

	ServiceInstance.taskManager = jobs.NewTaskManger(db)

	if ServiceInstance.taskManager.StartTaskManger() != nil {
		log.Fatal("error starting task manager")
	}

	SocketConn = domainSocket(socketDir)
	handleConnection(SocketConn)
}

func domainSocket(socketDir string) *net.UnixListener {

	var (
		conn     *net.UnixListener
		err      error
		unixAddr *net.UnixAddr
	)

	if err = os.MkdirAll(socketDir, 0755); err != nil {
		log.Fatal("error creating socket dir:", err)
	}

	socketFilePath := socketDir + "daemon.sock"

	syscall.Unlink(socketFilePath)

	if unixAddr, err = net.ResolveUnixAddr("unix", socketFilePath); err != nil {
		log.Fatal(err)
	}

	if conn, err = net.ListenUnix("unix", unixAddr); err != nil {
		log.Fatal("error creating domain socket:", err)
	}

	conn.SetUnlinkOnClose(true)

	return conn
}

func handleConnection(listener *net.UnixListener) {
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
		go treatMessage(c)
	}
}

func treatMessage(c net.Conn) {
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
			msg.WeekStatResponse, err = ServiceInstance.getWeekStat(msg.WeekStatRequest)

		case "appStat":
			msg.AppStatResponse, err = ServiceInstance.getAppStat(msg.AppStatRequest)

		case "dayStat":
			msg.DayStatResponse, err = ServiceInstance.getDayStat(msg.DayStatRequest)

		case "setCategory":
			msg.SetCategoryResponse, err = ServiceInstance.setAppCategory(msg.SetCategoryRequest)

		case "tasks":
			msg.ReminderAndLimitResponse, err = ServiceInstance.tasks()

		case "reminders":
			msg.ReminderAndLimitResponse, err = ServiceInstance.reminderTasks()

		case "appLimits":
			msg.ReminderAndLimitResponse, err = ServiceInstance.limitTasks()

		case "newReminder":
			msg.ReminderAndLimitResponse, err = ServiceInstance.addNewReminder(msg.ReminderAndLimitRequest)

		case "newAppLimit":
			msg.ReminderAndLimitResponse, err = ServiceInstance.addNewLimitApp(msg.ReminderAndLimitRequest)

		case "removeTask":
			msg.ReminderAndLimitResponse, err = ServiceInstance.removeTask(msg.ReminderAndLimitRequest)
		}

		if err != nil {
			msg.IsError, msg.Error = true, err.Error()
		}

		bytes, err := helperFuncs.EncodeJSON(msg)
		if err != nil {
			msg.Error = fmt.Errorf("error encoding response in serviceInstance: %w: %w", msg.Error, err).Error()
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
