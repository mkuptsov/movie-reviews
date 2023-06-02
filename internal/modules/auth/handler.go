package auth

import (
	"net/http"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/modules/echox"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	authService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Register(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.RegisterUserRequest](c)
	if err != nil {
		return err
	}

	user := &users.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     users.UserRole,
	}

	if err := h.authService.Register(c.Request().Context(), user, req.Password); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Login(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.LoginUserRequest](c)
	if err != nil {
		return err
	}

	accessToken, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	response := LoginResponse{
		AccessToken: accessToken,
	}
	return c.JSON(http.StatusOK, response)
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
