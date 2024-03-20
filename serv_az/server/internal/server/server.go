package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"server/internal/apperrors"
	"server/internal/config"
	"server/internal/database"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	"server/internal/log"
	"server/internal/repositories"
	"server/internal/services"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	repoLibrary repositories.RepoLibrary
	repoWords   repositories.RepoWords
	repoUsers   repositories.RepoUsers
	repoBackUp  repositories.BackUpCopyRepo
	router      Router
	logger      *logrus.Logger
	config      *config.Config
	blacklist   *blacklist
	tmpls       map[string]*template.Template
}

func NewServer(repoLibrary repositories.RepoLibrary, repoWords repositories.RepoWords, repoUsers repositories.RepoUsers, repoBackUp repositories.BackUpCopyRepo,
	logger *logrus.Logger, config *config.Config) *server {
	return &server{
		repoLibrary: repoLibrary,
		repoWords:   repoWords,
		repoUsers:   repoUsers,
		repoBackUp:  repoBackUp,
		router:      &router{mux: mux.NewRouter()},
		logger:      logger, config: config,
	}
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHttp(w, r)
}

func (srv *server) initializeRoutes() {
	srv.initializeTemplates()
	srv.logger.Info("server INIT")
	srv.router.InitImages("/images/", "templates/images")
	//------------HOME----translate
	srv.router.Get("/", srv.contextExpire(srv.homeHandler()))
	srv.router.GetPost("/translate", srv.contextExpire(srv.getTranslationHandler()))
	//---------------user-CRUD---------------------------
	srv.router.GetPost("/registration", srv.contextExpire(srv.createUserHandler()))
	srv.router.GetPost("/login", srv.contextExpire(srv.loginHandler()))
	blackList := newBlacklist()
	srv.blacklist = blackList ///user-update-password
	srv.router.GetPost("/user-update", srv.jwtAuthentication(srv.updateUserHandler()))
	srv.router.GetPost("/user-update-password", srv.jwtAuthentication(srv.updateUserPasswordHandler()))
	srv.router.GetPost("/logout", srv.contextExpire(srv.logoutHandler()))
	srv.router.Get("/user-info", srv.jwtAuthentication(srv.getUserByIdHandler()))
	//-------TESTS---LEARN--------------------
	srv.router.GetPost("/test", srv.jwtAuthentication(srv.testHandler()))
	srv.router.GetPost("/learn", srv.jwtAuthentication(srv.learnHandler()))
	//-------TESTS--------thematic test----------------------
	srv.router.GetPost("/thematic/{theme:[0-9a-zA-Z_-]+}", srv.jwtAuthentication(srv.testUniversalHandler()))
	srv.router.GetPost("/test-thematic", srv.jwtAuthentication(srv.themesHandler()))
	//-------Update LIBRARY-------------------
	srv.router.GetPost("/info-users", srv.jwtAuthentication(srv.getAllUsersHandler()))
	srv.router.GetPost("/library-update", srv.jwtAuthentication(srv.updateLibraryHandler()))
	srv.router.Get("/library-download", srv.jwtAuthentication(srv.downloadHandler()))
}

func Run() {
	logger, err := log.NewLogAndSetLevel("info")
	if err != nil {
		logger.Fatal(err)
	}

	cfg, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	if err = log.SetLevel(logger, cfg.Postgres.LogLevel); err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()
	psglDB := database.NewPostgresDB()
	db, err := psglDB.SetupDatabase(ctx, cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		if err := psglDB.PingEveryMinuts(ctx, 30, logger); err != nil {
			logger.Error(err)
		}
	}()

	if !db.Migrator().HasTable(&models.Library{}) {
		err = db.AutoMigrate(&models.Library{})
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Migration library OK")
	}

	if !db.Migrator().HasTable(&models.User{}) {
		err = db.AutoMigrate(&models.User{})
		if err != nil {
			logger.Fatal(err)
		}

		repoUser := repositories.NewRepoUsers(db, logger)
		repoLibrary := repositories.NewRepoLibrary(db, logger)
		repoWords := repositories.NewWords(db, logger)
		userService := services.NewUserService(repoWords, repoUser, repoLibrary, logger)
		adminUserReq := requests.CreateUserRequest{
			Email:    "admin@admin.admin",
			Name:     "mainName",
			LastName: "mainLastName",
			Password: cfg.Postgres.Password,
			Role:     "admin",
		}
		_, err = userService.CreateUser(ctx, &adminUserReq)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Migration Users OK")
		logger.Info("Admin created")
	}

	repoLibrary := repositories.NewRepoLibrary(db, logger)
	err = repoLibrary.InitWordsMap()
	if err != nil {
		logger.Fatal(err)
	}

	repoWords := repositories.NewWords(db, logger)
	repoUser := repositories.NewRepoUsers(db, logger)
	repoBackup := repositories.NewBackUpCopyRepo("save_copy/library.xlsx", logger)

	srv := NewServer(repoLibrary, repoWords, repoUser, repoBackup, logger, cfg)

	srv.initializeRoutes()
	logger.Infof("Listening HTTP service on %s port", cfg.AppPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), srv)
	if err != nil {
		logger.Fatal(err)
	}

	/*
		err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	*/
}

func (srv *server) initializeTemplates() error {
	tmplsList := make(map[string]*template.Template)

	tmpl, err := template.ParseFiles("templates/translate.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[translate] = tmpl

	tmpl, err = template.ParseFiles("templates/home.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[home] = tmpl

	tmpl, err = template.ParseFiles("templates/registration.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[registration] = tmpl

	tmpl, err = template.ParseFiles("templates/authentification.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[authentification] = tmpl

	tmpl, err = template.ParseFiles("templates/logout.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[logout] = tmpl

	tmpl, err = template.ParseFiles("templates/user_info.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[userInfo] = tmpl

	tmpl, err = template.ParseFiles("templates/test.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[test] = tmpl

	tmpl, err = template.ParseFiles("templates/test_thematic_handler.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[testThematicHandler] = tmpl

	tmpl, err = template.ParseFiles("templates/learn.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[learn] = tmpl

	tmpl, err = template.ParseFiles("templates/err.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[errMes] = tmpl

	tmpl, err = template.ParseFiles("templates/registrate.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[registrate] = tmpl

	tmpl, err = template.ParseFiles("templates/test_thematic.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[testThematic] = tmpl

	tmpl, err = template.ParseFiles("templates/update_lib.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[updateLib] = tmpl

	tmpl, err = template.ParseFiles("templates/update_user.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[updateUser] = tmpl

	tmpl, err = template.ParseFiles("templates/update_user_password.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[updateUserPassword] = tmpl

	tmpl, err = template.ParseFiles("templates/users_info.html", header, footer)
	if err != nil {
		appErr := apperrors.InitializeTemplatesErr.AppendMessage(err)
		srv.logger.Error(appErr)
		return appErr
	}
	tmplsList[usersInfo] = tmpl

	srv.tmpls = tmplsList

	srv.logger.Info("Templates have been registered")
	return nil
}

// "update_lib.html"
