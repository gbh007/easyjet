package httpapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt Controller) project(c echo.Context) error {
	ctx := c.Request().Context()

	var req projectIDRequest

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = c.Validate(&req)
	if err != nil {
		return err
	}

	project, err := cnt.service.Project(ctx, req.ProjectID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, project)
}
