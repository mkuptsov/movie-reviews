package echox

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mkuptsov/movie-reviews/internal/apperrors"
	"github.com/mkuptsov/movie-reviews/internal/log"
)

type HTTPError struct {
	Message    string `json:"message"`
	IncidentID string `json:"incident_id,omitempty"`
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appError *apperrors.Error
	if !errors.As(err, &appError) {
		appError = apperrors.InternalWithoutStackTrace(err)
	}

	httpError := HTTPError{
		Message:    appError.SafeError(),
		IncidentID: appError.IncidentID,
	}

	logger := log.FromContext(c.Request().Context())

	if appError.Code == apperrors.InternalCode {
		logger.Error("server error",
			"message", err.Error(),
			"incidentId", appError.IncidentID,
			"stack trace", appError.StackTrace)
	} else {
		logger.Warn("client error",
			"message", err.Error())
	}

	if err = c.JSON(toHTTPStatus(appError.Code), httpError); err != nil {
		logger.Error("server error",
			"message", err.Error())
	}
}

func toHTTPStatus(code apperrors.Code) int {
	switch code {
	case apperrors.InternalCode:
		return http.StatusInternalServerError
	case apperrors.BadRequestCode:
		return http.StatusBadRequest
	case apperrors.NotFoundCode:
		return http.StatusNotFound
	case apperrors.AlreadyExistsCode, apperrors.VersionMismatchCode:
		return http.StatusConflict
	case apperrors.UnauthorizedCode:
		return http.StatusUnauthorized
	case apperrors.ForbiddenCode:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
