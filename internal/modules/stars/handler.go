package stars

import (
	"net/http"

	"golang.org/x/sync/singleflight"

	"github.com/labstack/echo/v4"
	"github.com/mkuptsov/movie-reviews/contracts"
	"github.com/mkuptsov/movie-reviews/internal/config"
	"github.com/mkuptsov/movie-reviews/internal/echox"
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

func (h *Handler) CreateStar(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateStarRequest](c)
	if err != nil {
		return err
	}

	star := &StarDetails{
		Star: Star{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			BirthDate: req.BirthDate,
			DeathDate: req.DeathDate,
		},
		MiddleName: req.MiddleName,
		BirthPlace: req.BirthPlace,
		Bio:        req.Bio,
	}

	err = h.Service.CreateStar(c.Request().Context(), star)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, star)
}

func (h *Handler) GetStarByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetStarByIDRequest](c)
	if err != nil {
		return err
	}

	star, err := h.Service.GetStarByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, star)
}

func (h *Handler) GetAll(c echo.Context) error {
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[contracts.GetStarsRequest](c)
		if err != nil {
			return nil, err
		}

		pagination.SetDefaults(&req.PaginatedRequest, h.PaginationConfig)
		offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

		stars, total, err := h.Service.GetAllPaginated(c.Request().Context(), req.MovieID, offset, limit)
		if err != nil {
			return nil, err
		}

		return pagination.Response(&req.PaginatedRequest, total, stars), nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateStar(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateStarRequest](c)
	if err != nil {
		return err
	}

	star := &StarDetails{
		Star: Star{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			BirthDate: req.BirthDate,
			DeathDate: req.DeathDate,
		},
		MiddleName: req.MiddleName,
		BirthPlace: req.BirthPlace,
		Bio:        req.Bio,
	}
	id := req.ID

	err = h.Service.UpdateStar(c.Request().Context(), id, star)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteStar(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteStarRequest](c)
	if err != nil {
		return err
	}

	err = h.Service.DeleteStar(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
