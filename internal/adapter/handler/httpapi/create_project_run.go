package httpapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt Controller) createProjectRun(c echo.Context) error {
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

	id, err := cnt.service.RunProject(ctx, req.ProjectID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id": id,
	})
}
