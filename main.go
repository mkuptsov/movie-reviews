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
	"github.com/cloudmachinery/movie-reviews/internal/modules/auth"
	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/cloudmachinery/movie-reviews/internal/validation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "getting config: ")

	validation.SetupValidators()

	fmt.Printf("\nstarted with config:\n%+v\n", *cfg)

	db, err := getDb(context.Background(), cfg.DbUrl)
	failOnError(err, "getting db: ")
	defer db.Close()

	jwtService := jwt.NewService(cfg.Jwt.Secret, cfg.Jwt.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(usersModule.Service, jwtService)

	e := echo.New()
	api := e.Group("/api")

	authMiddleware := jwt.NewAuthMidlleware(cfg.Jwt.Secret)

	api.POST("/auth/register", authModule.Handler.Register)
	api.POST("/auth/login", authModule.Handler.Login)

	api.GET("/users", usersModule.Handler.GetUsers)
	api.GET("/users/:userId", usersModule.Handler.GetUserById)
	api.DELETE("/users/:userId", usersModule.Handler.Delete, authMiddleware, auth.Self)
	api.PUT("/users/:userId", usersModule.Handler.Update, authMiddleware, auth.Self)

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
}

func getDb(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connaction to db: %w", err)
	}

	return db, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
