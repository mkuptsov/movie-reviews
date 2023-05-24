package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *Service
}

type UpdateRequest struct {
	UserId int     `param:"userId" validate:"nonzero"`
	Bio    *string `json:"bio"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetUsers(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func (h *Handler) GetUserById(c echo.Context) error {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := h.Service.GetUserById(c.Request().Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.JSON(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c echo.Context) error {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.Service.Delete(c.Request().Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.JSON(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "user deleted")
}

func (h *Handler) Update(c echo.Context) error {
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return nil
	}

	err := h.Service.Update(c.Request().Context(), req.UserId, *req.Bio)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.JSON(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "user updated")
}
