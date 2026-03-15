package httpapi

import (
	"net/http"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/labstack/echo/v4"
)

func (cnt Controller) createProject(c echo.Context) error {
	ctx := c.Request().Context()

	var req entity.Project

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = c.Validate(&req)
	if err != nil {
		return err
	}

	id, err := cnt.service.CreateProject(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id": id,
	})
}
