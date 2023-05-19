package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (*Handler) GetUsers(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}
