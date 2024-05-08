package service

import (
	"LiScreMon/daemon/internal/database/repository"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
)

var (
	ServiceInstance Service
	SocketConn      *net.UnixListener
)

func StartService(socketDir string, db *repository.BadgerDBStore) {
	ServiceInstance.store = db
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
				fmt.Println("error accepting connection...uh we muuve:", err)
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
			msg Message
			err error
		)

		if err = gob.NewDecoder(c).Decode(&msg); err != nil {
			log.Println("error reading message:", err)
			continue
		}

		switch msg.Endpoint {
		case "startConnection":
			msg = Message{StringDataResponse: `hELLo.., this is the DaemonService speaking, your connection is established.`}

		case "closeConnection":
			fmt.Println("we got a close connection message")
			c.Close()
			return

		case "weekStat":
			weekStat := ServiceInstance.weekStat(msg)
			msg.WeekStatResponse = weekStat
		}

		bytes, err := msg.encode()
		if err != nil {
			log.Println("error encoding response:", err)
			continue
		}
		_, err = c.Write(bytes)
		if err != nil {
			log.Println("error encoding response:", err)
			continue
		}
	}
}
