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
	"github.com/cloudmachinery/movie-reviews/internal/modules/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/modules/auth"
	"github.com/cloudmachinery/movie-reviews/internal/modules/echox"
	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/cloudmachinery/movie-reviews/internal/validation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/validator.v2"
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "getting config: ")

	validation.SetupValidators()

	db, err := getDb(context.Background(), cfg.DbUrl)
	failOnError(err, "getting db: ")
	defer db.Close()

	jwtService := jwt.NewService(cfg.Jwt.Secret, cfg.Jwt.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(usersModule.Service, jwtService)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = RegisterAdmin(ctx, authModule, cfg.Admin)
	if apperrors.Is(err, apperrors.InternalCode) {
		log.Fatal(err)
	}

	e := echo.New()
	e.HTTPErrorHandler = echox.ErrorHandler
	e.Use(middleware.Recover())

	authMiddleware := jwt.NewAuthMidlleware(cfg.Jwt.Secret)
	api := e.Group("/api")
	api.Use(authMiddleware)

	api.POST("/auth/register", authModule.Handler.Register)
	api.POST("/auth/login", authModule.Handler.Login)

	api.GET("/users/:userId", usersModule.Handler.GetUserById)
	api.DELETE("/users/:userId", usersModule.Handler.Delete, auth.Self)
	api.PUT("/users/:userId", usersModule.Handler.Update, auth.Self)
	api.PUT("/users/:userId/role/:roleName", usersModule.Handler.UpdateUserRole, auth.Admin)

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

func RegisterAdmin(ctx context.Context, authModule *auth.Module, cfg config.AdminConfig) error {
	req := auth.RegisterRequest{
		Email:    cfg.Email,
		Username: cfg.Username,
		Pasword:  cfg.Password,
		Role:     users.AdminRole,
	}

	if cfg.Email == "" && cfg.Username == "" && cfg.Password == "" {
		return nil
	}
	err := validator.Validate(req)
	if err != nil {
		return apperrors.Internal(err)
	}

	return authModule.Service.Register(ctx, &users.User{
		Email:    req.Email,
		Username: req.Username,
		Role:     users.AdminRole,
	}, req.Pasword)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
