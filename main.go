package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/modules/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/modules/auth"
	"github.com/cloudmachinery/movie-reviews/internal/modules/echox"
	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/cloudmachinery/movie-reviews/internal/modules/log"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/cloudmachinery/movie-reviews/internal/validation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
	"gopkg.in/validator.v2"
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "getting config: ")

	validation.SetupValidators()

	logger, err := log.SetupLogger(cfg.Local, cfg.LogLevel)
	failOnError(err, "setup logger: ")
	slog.SetDefault(logger)

	db, err := getDb(context.Background(), cfg.DbUrl)
	failOnError(err, "getting db: ")
	defer db.Close()
	slog.Info("db connected")

	jwtService := jwt.NewService(cfg.Jwt.Secret, cfg.Jwt.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(usersModule.Service, jwtService)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = registerAdmin(ctx, authModule, cfg.Admin)
	failOnError(err, "create admin")

	e := echo.New()
	e.HTTPErrorHandler = echox.ErrorHandler
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())

	authMiddleware := jwt.NewAuthMidlleware(cfg.Jwt.Secret)
	api := e.Group("/api")
	api.Use(authMiddleware)
	api.Use(echox.Logger)

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

		slog.Info("graceful shutdown started")
		err := e.Shutdown(ctx)
		if err != nil {
			slog.Error("server shutdown", "error", err)
		}
	}()

	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
	}

	slog.Info("server stopped", "message", err)
}

func getDb(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connection to db: %w", err)
	}

	return db, nil
}

func registerAdmin(ctx context.Context, authModule *auth.Module, cfg config.AdminConfig) error {
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

	err = authModule.Service.Register(ctx, &users.User{
		Email:    req.Email,
		Username: req.Username,
		Role:     users.AdminRole,
	}, req.Pasword)

	if apperrors.Is(err, apperrors.InternalCode) {
		return err
	}
	if err == nil {
		slog.Info("admin created",
			"admin email", req.Email)
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}
