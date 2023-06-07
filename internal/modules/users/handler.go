package users

import (
	"net/http"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetUserByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetUserByIDRequest](c)
	if err != nil {
		return err
	}
	user, err := h.Service.GetUserByID(c.Request().Context(), req.UserID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserByUserName(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetUserByUserNameRequest](c)
	if err != nil {
		return err
	}

	user, err := h.Service.GetUserByUserName(c.Request().Context(), req.UserName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteUserRequest](c)
	if err != nil {
		return err
	}

	err = h.Service.DeleteUser(c.Request().Context(), req.UserID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "user deleted")
}

func (h *Handler) Update(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateUserRequest](c)
	if err != nil {
		return err
	}

	err = h.Service.Update(c.Request().Context(), req.UserID, *req.Bio)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) SetUserRole(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.SetUserRoleRequest](c)
	if err != nil {
		return err
	}

	err = h.Service.SetUserRole(c.Request().Context(), req.UserID, req.Role)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
