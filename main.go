package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/server"
	"golang.org/x/exp/slog"
)

const (
	gracefulTimeout = 10 * time.Second
)

func main() {
	fmt.Println(os.Getpid())

	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	srv, err := server.New(context.Background(), cfg)
	failOnError(err, "create server")

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
		defer cancel()

		if err = srv.Shutdown(ctx); err != nil {
			slog.Error("server shutdown", "error", err)
		}
	}()

	if err = srv.Start(); err != http.ErrServerClosed {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}
