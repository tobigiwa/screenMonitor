package service

import (
	"LiScreMon/daemon/internal/database/repository"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

type Service struct {
	service repository.IRepository
}

var ServiceInstance Service

func StartService(ctx context.Context, path string, db *repository.BadgerDBStore) {
	ServiceInstance.service = db
	handleConnection(ctx, domainSocket(path))
}

func domainSocket(configDir string) *net.UnixListener {

	var (
		conn     *net.UnixListener
		err      error
		unixAddr *net.UnixAddr
	)

	// dirPath := configDir + "/socket/"
	// dirPath := configDir + "/js/"

	// if err = os.MkdirAll(dirPath, 0755); err != nil {
	// 	log.Fatal("error creating socket dir:", err)
	// }

	// socketFilePath := dirPath + "s.sock"
	// _, err = os.Stat(socketFilePath)
	// fileNotExist := errors.Is(err, fs.ErrNotExist)

	// if err != nil && !fileNotExist {
	// 	log.Fatalln("error checking for existing socket file:", err)
	// }

	h, _ := os.UserHomeDir()
	socketFilePath := h + "/mod/s.sock"

	if err := syscall.Unlink(socketFilePath); err != nil {
		fmt.Println("error unlinking", err)
	}
	fmt.Println("unlinking successful")
	

	if _, err = os.Create(socketFilePath); err != nil {
		log.Fatal("error creating socket file:", err)
	}

	if unixAddr, err = net.ResolveUnixAddr("unix", socketFilePath); err != nil {
		fmt.Println("it came from here", err)
		log.Fatal(err)
	}

	if conn, err = net.ListenUnix("unix", unixAddr); err != nil {
		log.Fatal("error creating domain socket:", err)
	}

	conn.SetUnlinkOnClose(true)

	fmt.Println("socket is live")
	return conn
}

func handleConnection(ctx context.Context, conn *net.UnixListener) {
	for {
		_ = conn.SetDeadline(time.Now().Add(time.Second))
		c, err := conn.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() { // Ignore timeout errors
				continue
			}
			log.Fatal("error accepting connection:", err)
			break
		}

		fmt.Println("Connection accepted")

		go treatMessage(ctx, c)
	}
}

// func handleConnection2(ctx context.Context, conn net.Listener) {

// 	// connectionCh := make(chan net.Conn)

// 	go func() {
// 		for {
// 			c, err := conn.Accept()
// 			if err != nil {
// 				log.Fatal("error accepting connection:", err)
// 			}
// 			fmt.Println("Connection accepted")
// 			// connectionCh <- c
// 			wg.Add(1)

// 			go func() {
// 				defer wg.Done()

// 				treatMessage(c)
// 			}()
// 		}
// 	}()
// 	wg.Wait()
// for {
// 	select {
// 	case <-ctx.Done():
// 		fmt.Println("Received kill signal, shutting down...")
// 		conn.Close()
// 		return
// 	case c := <-connectionCh:
// 		go treatMessage(c)
// 	}
// }
// }
