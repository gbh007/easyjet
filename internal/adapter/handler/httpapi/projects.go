package httpapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt Controller) projects(c echo.Context) error {
	ctx := c.Request().Context()

	projects, err := cnt.service.Projects(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"projects": projects,
	})
}
