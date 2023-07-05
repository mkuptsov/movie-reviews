package echox

import (
	"github.com/labstack/echo/v4"
	"github.com/mkuptsov/movie-reviews/internal/jwt"
	"github.com/mkuptsov/movie-reviews/internal/log"
	"golang.org/x/exp/slog"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestGroup := slog.Group("request",
			slog.String("method", c.Request().Method),
			slog.String("url", c.Request().URL.String()))
		attrs := []any{requestGroup}

		if claims := jwt.GetClaims(c); claims != nil {
			requesterGroup := slog.Group("requester",
				slog.Int("id", claims.UserID),
				slog.String("role", claims.Role))

			attrs = append(attrs, requesterGroup)
		}

		logger := slog.Default().With(attrs...)
		ctx := log.WithLogger(c.Request().Context(), logger)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
