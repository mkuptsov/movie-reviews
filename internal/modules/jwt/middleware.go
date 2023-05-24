package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	tokenContexKey = "token"
)

func NewAuthMidlleware(secret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey: tokenContexKey,
		SigningKey: []byte(secret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &AccessClaims{}
		},
	})
}

func GetClaims(c echo.Context) *AccessClaims {
	token := c.Get(tokenContexKey)
	if token == nil {
		// return echo.NewHTTPError(http.StatusUnauthorized)
		panic("attempt to get token before auth middleware")
	}
	t, ok := token.(*jwt.Token)
	if !ok {
		panic("invalid token type")
	}

	ac, ok := t.Claims.(*AccessClaims)
	if !ok {
		panic("invalid claims type")
	}

	return ac
}
