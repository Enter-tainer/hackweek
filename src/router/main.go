package router

import "github.com/labstack/echo"

func InitRouter(e *echo.Group) {
	userGroup := e.Group("/user")
	postGroup := e.Group("/post")
	initUserGroup(userGroup)
	initPostGroup(postGroup)
}
