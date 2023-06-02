package users

import (
	"net/http"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/modules/echox"
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

func (h *Handler) GetUserById(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetUserByIdRequest](c)
	if err != nil {
		return err
	}
	user, err := h.Service.GetUserById(c.Request().Context(), req.UserId)
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

	err = h.Service.DeleteUser(c.Request().Context(), req.UserId)
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

	err = h.Service.Update(c.Request().Context(), req.UserId, *req.Bio)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) UpdateUserRole(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateUserRoleRequest](c)
	if err != nil {
		return err
	}

	err = h.Service.UpdateUserRole(c.Request().Context(), req.UserId, req.Role)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
