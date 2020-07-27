package middleware

import (
	"tree-hole/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWT([]byte(config.Config.JWT.Secret))
}
