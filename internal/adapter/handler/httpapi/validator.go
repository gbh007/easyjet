package httpapi

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type vldr struct {
	validator *validator.Validate
}

func (v vldr) Validate(i any) error {
	err := v.validator.Struct(i)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
