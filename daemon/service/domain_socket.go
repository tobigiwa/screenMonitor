package service

import (
	"encoding/gob"
	"log"
	"net"
	"os"
)

func StartService(path string, killSignal chan os.Signal) {
	handleConnection(domainSocket(path), killSignal)
}

func domainSocket(configDir string) *net.UnixListener {

	var (
		conn *net.UnixListener
		err  error
	)

	dirPath := configDir + "/socket/"
	if err = os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatal("error creating socket dir:", err)
	}

	if conn, err = net.ListenUnix("unix", &net.UnixAddr{Name: dirPath + "liScreMon.sock"}); err != nil {
		log.Fatal("error creating domain socket:", err)
	}

	return conn
}

func handleConnection(conn *net.UnixListener, killSignal chan os.Signal) {
	for {
		select {
		case <-killSignal:
			conn.Close()
		default:
			c, err := conn.Accept()
			if err != nil {
				log.Fatal("error accepting connection:", err)
			}
			go handleRequest(c)
		}
	}
}

func handleRequest(c net.Conn) {
	defer c.Close()

	var (
		msg Message
		err error
	)
	if err = gob.NewDecoder(c).Decode(&msg); err != nil {
		log.Println("error reading message:", err)
		return
	}

	switch msg.Type {
	case "hello":

	}

}

type Message struct {
	Type string
	Body string
}
