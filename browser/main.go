package main

import (
	"browser/webserver"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {

	opts := slog.HandlerOptions{
		AddSource: true,
	}

	textLogger := slog.NewTextHandler(os.Stdout, &opts)
	logger := slog.New(textLogger)
	slog.SetDefault(logger)

	app, err := webserver.NewApp(logger)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Fatalln("daemon service is not running", err)
		}
		log.Fatalln("error creating app:", err)
	}

	err = app.CheckDaemonService()
	if err != nil {
		log.Fatalln("error connectig to daemon service:", err)
	}

	log.Println("Connected to daemon service")

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.Routes(),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func(done chan os.Signal) {
		<-done
		close(done)

		if err := app.CloseDaemonConnection(); err != nil {
			fmt.Println("error closing socket connection with daemon, error:", err)
		}

		if err := server.Shutdown(context.TODO()); err != nil {
			log.Fatalf("Graceful server shutdown Failed:%+v\n", err)
		}
	}(done)

	fmt.Println("Server is running on port http://127.0.0.1:8080/screentime")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("Server error:", err)
	}
	fmt.Println()
	fmt.Println("SERVER STOPPED GRACEFULLY")
}
