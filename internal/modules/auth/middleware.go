package auth

import (
	"github.com/cloudmachinery/movie-reviews/internal/modules/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/labstack/echo/v4"
)

func Self(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims == nil {
			return apperrors.Unauthorized("invalid or missing token")
		}
		userID := c.Param("userId")
		if claims.Role == users.AdminRole || claims.Subject == userID {
			return next(c)
		}

		return apperrors.Forbidden("insufficient permissions")
	}
}

func Editor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims == nil {
			return apperrors.Unauthorized("invalid or missing token")
		}
		if claims.Role == users.AdminRole || claims.Role == users.EditorRole {
			return next(c)
		}

		return apperrors.Forbidden("insufficient permissions")
	}
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims == nil {
			return apperrors.Unauthorized("invalid or missing token")
		}
		if claims.Role == users.AdminRole {
			return next(c)
		}

		return apperrors.Forbidden("insufficient permissions")
	}
}
