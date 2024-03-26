package router

import (
	"server/internal/infrastructure/middleware"
	"server/internal/interface/controller"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
)

func NewRouter(e *echo.Echo, srv controller.AppController, secretKey string) *echo.Echo {
	//e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
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
	e.POST("/login", func(context echo.Context) error { return srv.HandlerController.LoginHandler(context) })
	e.GET("/login", func(context echo.Context) error { return srv.HandlerController.LoginHandler(context) })

	e.POST("/user-restore-password", func(context echo.Context) error { return srv.HandlerController.RestoreUserPasswordHandler(context) })
	e.GET("/user-restore-password", func(context echo.Context) error { return srv.HandlerController.RestoreUserPasswordHandler(context) })
	blackList := middleware.NewBlacklist()
	e.POST("/logout", srv.HandlerController.LogoutHandler(blackList))
	e.GET("/logout", srv.HandlerController.LogoutHandler(blackList))
	//---------------JWT-------------------------
	jwtConfig := middleware.JWTMiddlewareConfig{SecretKey: secretKey}

	e.GET("/user-info", srv.HandlerController.GetUserByIdHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.GET("/user-update", srv.HandlerController.UpdateUserHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.POST("/user-update", srv.HandlerController.UpdateUserHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.GET("/user-update-password", srv.HandlerController.UpdateUserPasswordHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.POST("/user-update-password", srv.HandlerController.UpdateUserPasswordHandler, middleware.JWTAuthentication(&jwtConfig, blackList))

	//-------Update LIBRARY-------------------
	e.GET("/info-users", srv.HandlerController.GetAllUsersHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.POST("/library-update", srv.HandlerController.UpdateLibraryHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.GET("/library-update", srv.HandlerController.UpdateLibraryHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.GET("/library-download", srv.HandlerController.DownloadHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	//-------TESTS---LEARN--------------------
	e.POST("/test", srv.HandlerController.TestHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.GET("/test", srv.HandlerController.TestHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.POST("/learn", srv.HandlerController.LearnHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.GET("/learn", srv.HandlerController.LearnHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	//-------TESTS--------thematic test----------------------
	e.POST("/thematic", srv.HandlerController.TestUniversalHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.GET("/thematic", srv.HandlerController.TestUniversalHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	//e.POST("/test-thematic", srv.HandlerController.ThemesHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	e.GET("/test-thematic", srv.HandlerController.ThemesHandler, middleware.JWTAuthentication(&jwtConfig, blackList))
	return e
}
