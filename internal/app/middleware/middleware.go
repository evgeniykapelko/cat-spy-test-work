package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Logger() echo.MiddlewareFunc {
	return middleware.Logger()
}

func Recover() echo.MiddlewareFunc {
	return middleware.Recover()
}
