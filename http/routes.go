package http

import (
	"net/http"

	"mckp/roberts-concordance/data"

	"github.com/labstack/echo/v4"
)

func HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "Healthy")
}

func ReadinessCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "Ready")
}

func GetBible(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, data.GetText())
}

func GetBooksOfBible(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, data.GetBooks())
}
