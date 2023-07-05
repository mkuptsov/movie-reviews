package reviews

import (
	"errors"
	"net/http"

	"golang.org/x/sync/singleflight"

	"github.com/labstack/echo/v4"
	"github.com/mkuptsov/movie-reviews/contracts"
	"github.com/mkuptsov/movie-reviews/internal/apperrors"
	"github.com/mkuptsov/movie-reviews/internal/config"
	"github.com/mkuptsov/movie-reviews/internal/echox"
	"github.com/mkuptsov/movie-reviews/internal/pagination"
)

type Handler struct {
	service          *Service
	paginationConfig config.PaginationConfig
	reqGroup         singleflight.Group
}

func NewHandler(service *Service, paginationConfig config.PaginationConfig) *Handler {
	return &Handler{
		service:          service,
		paginationConfig: paginationConfig,
	}
}

func (h *Handler) GetAll(c echo.Context) error {
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[contracts.GetReviewsRequest](c)
		if err != nil {
			return nil, err
		}

		if req.MovieID == nil && req.UserID == nil {
			return nil, apperrors.BadRequest(errors.New("either movie_id or user_id must be provided"))
		}

		pagination.SetDefaults(&req.PaginatedRequest, h.paginationConfig)
		offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

		reviews, total, err := h.service.GetPaginated(c.Request().Context(), req.MovieID, req.UserID, offset, limit)
		if err != nil {
			return nil, err
		}
		return pagination.Response(&req.PaginatedRequest, total, reviews), nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Get(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetReviewRequest](c)
	if err != nil {
		return err
	}

	review, err := h.service.GetByID(c.Request().Context(), req.ReviewID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, review)
}

func (h *Handler) Create(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateReviewRequest](c)
	if err != nil {
		return err
	}

	review := &Review{
		MovieID: req.MovieID,
		UserID:  req.UserID,
		Rating:  req.Rating,
		Title:   req.Title,
		Content: req.Content,
	}

	err = h.service.Create(c.Request().Context(), review)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, review)
}

func (h *Handler) Update(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateReviewRequest](c)
	if err != nil {
		return err
	}

	if err = h.service.Update(c.Request().Context(), req.ReviewID, req.UserID, req.Title, req.Content, req.Rating); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) Delete(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteReviewRequest](c)
	if err != nil {
		return err
	}

	if err = h.service.Delete(c.Request().Context(), req.ReviewID, req.UserID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
