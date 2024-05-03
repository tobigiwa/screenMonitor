package service

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/fs"
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

	socketFilePath := dirPath + "liScreMon.sock"

	if _, err = os.Stat(socketFilePath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatal("error checking for existing socket file:", err)
	}

	if err = os.Remove(socketFilePath); err != nil {
		log.Fatal("error removing existing socket file:", err)
	}

	if conn, err = net.ListenUnix("unix", &net.UnixAddr{Name: socketFilePath}); err != nil {
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
			go treatMessage(c)
		}
	}
}

func treatMessage(c net.Conn) {
	defer c.Close()

	var (
		msg Message
		err error
	)

	if err = gob.NewDecoder(c).Decode(&msg); err != nil {
		log.Println("error reading message:", err)
		return
	}

	fmt.Printf("\n\n%+v\n\n", msg)

	c.Write([]byte("hello from the daemon"))
}

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
