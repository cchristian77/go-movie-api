package utils

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-movie-api/token"
	"go-movie-api/utils/helper"
	"go-movie-api/utils/response"
	"net/http"
)

// ErrorHandler returns JSON including status code and error message if error occurs
func ErrorHandler(err error, ec echo.Context) {
	var statusCode int
	var errorMsg string

	// Get status code and error message from if error is HTTP error type
	httpError, ok := err.(*echo.HTTPError)
	if ok {
		statusCode = httpError.Code
		errorMsg = fmt.Sprintf("%s", httpError.Message)
	} else {
		statusCode = getStatusCode(err)
		errorMsg = err.Error()
	}

	// record error to log
	Logger.Error(errorMsg)

	// Return JSON with status code and error message
	if !ec.Response().Committed {
		ec.JSON(statusCode, response.Error{Message: errorMsg, Status: statusCode})
	}
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case helper.InternalServerErr:
		return http.StatusInternalServerError
	case helper.NotFoundErr:
		return http.StatusNotFound
	case helper.ConflictErr:
		return http.StatusConflict
	case helper.BadParamInputErr, helper.IncorrectCredentialErr:
		return http.StatusBadRequest
	case helper.ForbiddenErr:
		return http.StatusForbidden
	case helper.UnauthorizedErr, token.InvalidTokenErr, token.ExpiredTokenErr:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
