package service

import (
	db "LiScreMon/cli/daemon/internal/database"
	"encoding/json"
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
		)

		if err = json.NewDecoder(c).Decode(&msg); err != nil {
			fmt.Println("error reading message:", err)
			if errors.Is(err, io.EOF) {
				fmt.Println("client connection closed")
				c.Close()
				return
			}
			continue
		}

		switch msg.Endpoint {
		case "startConnection":
			msg = types.Message{StringDataResponse: `hELLo.., this is the DaemonService speaking, your connection is established.`}

		case "closeConnection":
			fmt.Println("we got a close connection message")
			c.Close()
			return

		case "weekStat":
			weekStat := ServiceInstance.getWeekStat(msg)
			msg.WeekStatResponse = weekStat
		}

		bytes, err := helperFuncs.Encode(msg)
		if err != nil {
			fmt.Println("error encoding response:", err)
			continue
		}
		_, err = c.Write(bytes)
		if err != nil {
			fmt.Println("error encoding response:", err)
			continue
		}
	}
}
