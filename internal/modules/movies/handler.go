package movies

import (
	"net/http"

	"golang.org/x/sync/singleflight"

	"github.com/labstack/echo/v4"
	"github.com/mkuptsov/movie-reviews/contracts"
	"github.com/mkuptsov/movie-reviews/internal/config"
	"github.com/mkuptsov/movie-reviews/internal/echox"
	"github.com/mkuptsov/movie-reviews/internal/modules/genres"
	"github.com/mkuptsov/movie-reviews/internal/modules/stars"
	"github.com/mkuptsov/movie-reviews/internal/pagination"
)

type Handler struct {
	Service          *Service
	PaginationConfig config.PaginationConfig
	reqGroup         singleflight.Group
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

	for _, item := range req.Genres {
		movie.Genres = append(movie.Genres, &genres.Genre{ID: item})
	}

	for _, item := range req.Cast {
		movie.Cast = append(movie.Cast, &stars.MovieCredit{
			Star: stars.Star{
				ID: item.StarID,
			},
			Role:    item.Role,
			Details: item.Details,
		})
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
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[contracts.GetMoviesRequest](c)
		if err != nil {
			return nil, err
		}

		pagination.SetDefaults(&req.PaginatedRequest, h.PaginationConfig)
		offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

		movies, total, err := h.Service.GetAllPaginated(c.Request().Context(), req.StarID, req.SearchTerm, req.SortByRating, offset, limit)
		if err != nil {
			return nil, err
		}
		return pagination.Response(&req.PaginatedRequest, total, movies), nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
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
		Version:     req.Version,
	}
	id := req.ID

	for _, item := range req.Genres {
		movie.Genres = append(movie.Genres, &genres.Genre{ID: item})
	}

	for _, item := range req.Cast {
		movie.Cast = append(movie.Cast, &stars.MovieCredit{
			Star: stars.Star{
				ID: item.StarID,
			},
			Role:    item.Role,
			Details: item.Details,
		})
	}

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
