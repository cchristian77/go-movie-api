package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type RequestValidator struct {
	Validator *validator.Validate
}

func (v *RequestValidator) Validate(request any) error {
	if err := v.Validator.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
