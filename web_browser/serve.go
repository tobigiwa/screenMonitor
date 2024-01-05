package browser

import (
	backend "browser/web_server"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

// Web launches a webserver for LiScreMon to be
// displayed on your local machine browser
func Web() error {

	// logging services
	opts := slog.HandlerOptions{
		AddSource: true,
	}

	jsonLogger := slog.NewTextHandler(os.Stdout, &opts)
	logger := slog.New(jsonLogger)
	slog.SetDefault(logger)

	app := backend.NewApp(logger)

	// if !logging {
	// 	log.SetOutput(io.Discard)

	// }

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.Routes(),
	}

	fmt.Println("Server is running on port http://127.0.0.1:8080/screentime")

	if err := server.ListenAndServe(); err != nil {
		return err
	}
	fmt.Println("SERVING DONE")

	return nil
}
