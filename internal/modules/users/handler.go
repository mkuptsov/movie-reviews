package users

import (
	"net/http"
	"strconv"

	"github.com/cloudmachinery/movie-reviews/internal/modules/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/modules/echox"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *Service
}

type UpdateRequest struct {
	UserId int     `param:"userId" validate:"nonzero"`
	Bio    *string `json:"bio"`
}

type UpdateUserRoleRequest struct {
	userId   int    `param:"userId" validate:"nonzero"`
	roleName string `param:"roleName" validate:"nonzero,role"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetUserById(c echo.Context) error {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return apperrors.BadRequestHidden(err, "userId must be a number")
	}
	user, err := h.Service.GetUserById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c echo.Context) error {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return apperrors.BadRequestHidden(err, "userId must be a number")
	}

	err = h.Service.Delete(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "user deleted")
}

func (h *Handler) Update(c echo.Context) error {
	req, err := echox.BindAndValidate[UpdateRequest](c)
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
	req, err := echox.BindAndValidate[UpdateUserRoleRequest](c)
	if err != nil {
		return err
	}

	err = h.Service.UpdateUserRole(c.Request().Context(), req.userId, req.roleName)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
