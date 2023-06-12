package stars

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
	req, err := echox.BindAndValidate[contracts.GetStarsRequest](c)
	if err != nil {
		return err
	}

	pagination.SetDefaults(&req.PaginatiedRequest, h.PaginationConfig)
	offset, limit := pagination.OffsetLimit(&req.PaginatiedRequest)

	stars, total, err := h.Service.GetAllPaginated(c.Request().Context(), offset, limit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.Response(&req.PaginatiedRequest, total, stars))
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
