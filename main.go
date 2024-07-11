package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
)

func main() {
	appMode := os.Getenv("APP_MODE")
	if strings.TrimSpace(appMode) == "" {
		appMode = "dev"
		log.Println("APP_MODE: (default) dev")
	}

	app := controllers.App{
		appMode: appMode,
	}

	// Gracefully Shutdown
	shutdown := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Println("Interrupt signal received. Shutting down.")
		err := app.Server.Shutdown(context.WithTimeout(context.Background(), 60))
		if err != nil {
			log.Println("HTTP Server shutdown error: \n", err.Error())
		}
		close(shutdown)
	}()

	app.Run()

	<-shutdown
}
