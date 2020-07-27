package main

import (
	"log"
	"tree-hole/config"
	"tree-hole/model"
	"tree-hole/router"
	"tree-hole/util"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	config.InitConfig()
	model.InitModel()
	util.InitUtil()

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Validator = &CustomValidator{
		validator: validator.New(),
	}

	apiGroup := e.Group("/api/v1")
	router.InitRouter(apiGroup)

	log.Fatal(e.Start(config.Config.App.Address))
}
