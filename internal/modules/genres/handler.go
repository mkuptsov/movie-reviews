package genres

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

func (h *Handler) GetGenres(c echo.Context) error {
	genres, err := h.Service.GetGenres(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, genres)
}

func (h *Handler) GetGenreByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetGenreByIDRequest](c)
	if err != nil {
		return err
	}

	genre, err := h.Service.GetGenreByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, genre)
}

func (h *Handler) CreateGenre(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateGenreRequest](c)
	if err != nil {
		return err
	}

	genre, err := h.Service.CreateGenre(c.Request().Context(), req.Name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, genre)
}

func (h *Handler) UpdateGenre(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateGenreRequest](c)
	if err != nil {
		return err
	}

	return h.Service.UpdateGenre(c.Request().Context(), req.ID, req.Name)
}

func (h *Handler) DeleteGenre(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteGenreRequest](c)
	if err != nil {
		return err
	}

	return h.Service.DeleteGenre(c.Request().Context(), req.ID)
}
