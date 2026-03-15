package httpapi

import (
	"net/http"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/labstack/echo/v4"
)

func (cnt Controller) updateProject(c echo.Context) error {
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

	err = cnt.service.UpdateProject(ctx, req)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
