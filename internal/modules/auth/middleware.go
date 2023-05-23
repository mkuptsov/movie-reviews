package auth

import (
	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/labstack/echo/v4"
)

func Self(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("userId")
		claims := jwt.GetClaims(c)
		if claims.Role == "admin" || claims.Subject == userID {
			return next(c)
		}

		return echo.ErrForbidden
	}
}

func Editor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims.Role == "admin" || claims.Role == "editor" {
			return next(c)
		}

		return echo.ErrForbidden
	}
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims.Role == "admin" {
			return next(c)
		}

		return echo.ErrForbidden
	}
}
