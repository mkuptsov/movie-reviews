package movies

import (
	"net/http"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"github.com/cloudmachinery/movie-reviews/internal/pagination"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service          *Service
	PaginationConfig config.PaginationConfig
}

func NewHandler(service *Service, cfg config.PaginationConfig) *Handler {
	return &Handler{
		Service:          service,
		PaginationConfig: cfg,
	}
}

func (h *Handler) CreateMovie(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateMovieRequest](c)
	if err != nil {
		return err
	}

	movie := &MovieDetails{
		Movie: Movie{
			Title:       req.Title,
			ReleaseDate: req.ReleaseDate,
		},
		Description: req.Description,
	}

	err = h.Service.CreateMovie(c.Request().Context(), movie)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, movie)
}

func (h *Handler) GetMovieByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetMovieByIDRequest](c)
	if err != nil {
		return err
	}

	movie, err := h.Service.GetMovieByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

func (h *Handler) GetAll(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetMoviesRequest](c)
	if err != nil {
		return err
	}

	pagination.SetDefaults(&req.PaginatiedRequest, h.PaginationConfig)
	offset, limit := pagination.OffsetLimit(&req.PaginatiedRequest)

	movies, total, err := h.Service.GetAllPaginated(c.Request().Context(), offset, limit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.Response(&req.PaginatiedRequest, total, movies))
}

func (h *Handler) UpdateMovie(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateMovieRequest](c)
	if err != nil {
		return err
	}

	movie := &MovieDetails{
		Movie: Movie{
			Title:       req.Title,
			ReleaseDate: req.ReleaseDate,
		},
		Description: req.Description,
	}
	id := req.ID

	err = h.Service.UpdateMovie(c.Request().Context(), id, movie)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteMovie(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteMovieRequest](c)
	if err != nil {
		return err
	}

	err = h.Service.DeleteMovie(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
