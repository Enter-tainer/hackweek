package router

import (
	"tree-hole/controller"
	middleware2 "tree-hole/middleware"

	"github.com/labstack/echo"
)

func initUserGroup(group *echo.Group) {
	group.GET("/token", controller.UserGetToken)
	group.POST("/info", controller.UserRegister)
	group.POST("/verify", controller.UserVerify)

	group.PUT("/info", controller.UserUpdateInfo, middleware2.JWTMiddleware())
	group.DELETE("", controller.UserDelete, middleware2.JWTMiddleware())
	group.GET("/info", controller.UserGetInfo, middleware2.JWTMiddleware())
}
