package http

import (
	"net/http"
	"strconv"

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

func GetSpecificBookOfBible(ctx echo.Context) error {
	bookName := ctx.QueryParam("book")

	return ctx.JSON(http.StatusOK, data.GetBookByName(bookName))
}

func GetVersesForBookOfBible(ctx echo.Context) error {
	bookName := ctx.Param("book")
	start, _ := strconv.Atoi(ctx.QueryParam("start"))
	end, _ := strconv.Atoi(ctx.QueryParam("end"))

	// start - 1 so that we use index and not the verse number.
	// i.e index 0 is verse 1
	return ctx.JSON(http.StatusOK, data.GetVerseFromBook(bookName, start-1, end))
}
