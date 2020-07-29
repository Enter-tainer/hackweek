package router

import (
	"tree-hole/controller"
	"tree-hole/middleware"

	"github.com/labstack/echo"
)

func initPostGroup(group *echo.Group) {
	group.GET("/", controller.PostGetAll, middleware.JWTMiddleware())
	group.GET("/:id", controller.PostGetWithId, middleware.JWTMiddleware())
	group.POST("/", controller.PostNew, middleware.JWTMiddleware())
	group.POST("/:id", controller.PostNewComment, middleware.JWTMiddleware())
}
