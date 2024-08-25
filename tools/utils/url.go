package utils

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func ConvertURL(c echo.Context, path string) string {
	host := c.Request().Host

	scheme := "http"

	h := c.Request().Header.Get("X-Forwarded-Proto")
	if h != "" {
		scheme = h
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}
