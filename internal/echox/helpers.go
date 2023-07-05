package echox

import (
	"github.com/labstack/echo/v4"
	"github.com/mkuptsov/movie-reviews/internal/apperrors"
	"gopkg.in/validator.v2"
)

func BindAndValidate[T any](c echo.Context) (*T, error) {
	req := new(T)
	if err := c.Bind(req); err != nil {
		return nil, apperrors.BadRequestHidden(err, "invalid or malformed request")
	}

	if err := validator.Validate(req); err != nil {
		return nil, apperrors.BadRequest(err)
	}
	return req, nil
}
