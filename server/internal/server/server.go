package server

import (
	"context"
	"fmt"
	"net/http"
	"server/internal/config"
	"server/internal/database"
	"server/internal/domain/models"
	"server/internal/log"
	"server/internal/repositories"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	repoLibrary repositories.RepoLibrary
	repoUsers   repositories.RepoUsers
	router      Router
	logger      *logrus.Logger
	config      *config.Config
	blacklist   *blacklist
}

func NewServer(repoLibrary repositories.RepoLibrary, repoUsers repositories.RepoUsers, logger *logrus.Logger, config *config.Config) *server {
	return &server{repoLibrary: repoLibrary, repoUsers: repoUsers, router: &router{mux: mux.NewRouter()}, logger: logger, config: config}
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHttp(w, r)
}

func (srv *server) initializeRoutes() {
	srv.logger.Info("server INIT")
	srv.router.Get("/library/translate", srv.contextExpire(srv.getTranslationHandler()))

	srv.router.Post("/users", srv.contextExpire(srv.createUserHandler()))
	srv.router.Post("/users/login", srv.contextExpire(srv.loginHandler()))

	blackList := newBlacklist()
	srv.blacklist = blackList
	srv.router.Post("/users/logout", srv.contextExpire(srv.logoutHandler()))
	srv.router.Get("/users/{user_id}", srv.jwtAuthentication(srv.getUserByIdHandler()))
	srv.router.Get("/user/words", srv.jwtAuthentication(srv.getWordsByUserIDAndLimitHandler()))
	srv.router.Put("/user/move-word-to-learned", srv.jwtAuthentication(srv.moveWordToLearnedHandler()))
	srv.router.Post("/user/add-word-to-learn", srv.jwtAuthentication(srv.addWordToLearnHandler()))
	srv.router.Get("/user/learn", srv.jwtAuthentication(srv.getLearnByUserIDAndLimitHandler()))
	srv.router.Delete("/user/learn", srv.jwtAuthentication(srv.deleteLearnByUserIDAndLearnIDHandler()))

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

	if !db.Migrator().HasTable(&models.Library{}) {
		err = db.AutoMigrate(&models.Library{})
		if err != nil {
			logger.Fatal(err)
		}

		repoBackup := repositories.NewBackUpCopyRepo("save_copy/library.json", "save_copy/library.txt", logger)
		words, err := repoBackup.GetAllFromBackUp()
		if err != nil {
			logger.Fatal(err)
		}

		repoLibrary := repositories.NewRepoLibrary(db, logger)
		err = repoLibrary.InsertWordsLibrary(ctx, words)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Migration success")
	}

	if !db.Migrator().HasTable(&models.User{}) {
		err = db.AutoMigrate(&models.User{})
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Migration success")
	}

	repoLibrary := repositories.NewRepoLibrary(db, logger)
	repoUser := repositories.NewRepoUsers(db, logger)
	srv := NewServer(repoLibrary, repoUser, logger, cfg)

	srv.initializeRoutes()
	logger.Infof("Listening HTTP service on %s port", cfg.AppPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), srv)
	if err != nil {
		logger.Fatal(err)
	}
}
