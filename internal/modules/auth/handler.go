package auth

import (
	"net/http"

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
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	user := &users.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := h.authService.Register(c.Request().Context(), user, req.Pasword); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	accessToken, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	}

	response := LoginResponse{
		AccessToken: accessToken,
	}
	return c.JSON(http.StatusOK, response)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pasword  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
