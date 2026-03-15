package httpapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt Controller) projectRun(c echo.Context) error {
	ctx := c.Request().Context()

	var req projectIDAndRunIDRequest

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = c.Validate(&req)
	if err != nil {
		return err
	}

	run, err := cnt.service.ProjectRun(ctx, req.RunID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, run)
}
