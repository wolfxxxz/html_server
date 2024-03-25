package router

import (
	"server/internal/interface/controller"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewRouter(e *echo.Echo, srv controller.AppController, secretKey string) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//-------init images ------------------
	e.Static("/images", "templates/images")
	//------------HOME----translate
	e.GET("/", func(context echo.Context) error { return srv.HandlerController.HomeHandler(context) })
	e.GET("/translate", func(context echo.Context) error { return srv.HandlerController.GetTranslationHandler(context) })
	e.POST("/translate", func(context echo.Context) error { return srv.HandlerController.GetTranslationHandler(context) })
	e.GET("/quick-answer", func(context echo.Context) error { return srv.HandlerController.QuickAnswerHandler(context) })
	//---------------user-CRUD----------------
	e.GET("/registration", func(context echo.Context) error { return srv.HandlerController.CreateUserHandler(context) })
	e.POST("/registration", func(context echo.Context) error { return srv.HandlerController.CreateUserHandler(context) })
	//srv.router.GetPost("/login", srv.contextExpire(srv.loginHandler()))

	return e
}
