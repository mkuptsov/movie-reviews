package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	cfg, err := config.NewConfig()
	failOnError(err, "getting config: ")

	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := e.Shutdown(ctx)
		if err != nil {
			fmt.Printf("shutdown server: %v", err)
		}
	}()

	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	fmt.Printf("server stopped: %v\n", err)
	// db.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
