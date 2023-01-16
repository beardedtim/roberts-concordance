package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var server *echo.Echo

func Create() {
	server = echo.New()
}

func attachMiddleware() {
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
}

func attachRoutes() {
	server.GET("/.well-known/healthcheck", HealthCheck)
	server.GET("/.well-known/readiness", ReadinessCheck)

	server.GET("/bible", GetBible)
	server.GET("/bible/books", GetBooksOfBible)
	server.GET("/bible/:book", GetSpecificBookOfBible)
	server.GET("/bible/:book/verses", GetVersesForBookOfBible)
}

func configure() {
	attachMiddleware()
	attachRoutes()
}

func Start(port int) error {
	configure()

	return server.Start(fmt.Sprintf(":%d", port))
}
