package middleware

import "github.com/labstack/echo/v4"

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(Logger())
	e.Use(Recover())
}
