package main

import (
	"context"
	"server/internal/config"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	"server/internal/domain/validator"
	"server/internal/infrastructure/datastore"
	"server/internal/infrastructure/email"
	"server/internal/infrastructure/router"
	"server/internal/infrastructure/webtemplate.go"
	"server/internal/interface/repository"
	"server/internal/log"
	"server/internal/registry"
	"server/internal/usercase/interactor"
	"time"

	"github.com/labstack/echo"
)

func main() {
	time.Sleep(3 * time.Second)
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
	psglDB := datastore.NewPostgresDB()
	db, err := psglDB.SetupDatabase(ctx, cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		if err := psglDB.PingEveryMinuts(ctx, 30, logger); err != nil {
			logger.Error(err)
		}
	}()

	sender := email.InitSender(cfg.Email.Email, cfg.Email.Key, cfg.Email.SMTP, cfg.Email.Port)
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

		repoUser := repository.NewUserRepository(db, logger)
		repoWords := repository.NewWordsRepository(db, logger)
		usInteractor := interactor.NewUserInteractor(repoUser, repoWords, sender)
		adminUserReq := requests.CreateUserRequest{
			Email:    "admin@admin.admin",
			Name:     "mainName",
			LastName: "mainLastName",
			Password: cfg.Postgres.Password,
			Role:     "admin",
		}

		_, err = usInteractor.CreateUser(ctx, &adminUserReq)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Migration Users OK")
		logger.Info("Admin created")
	}

	repoLibrary := repository.NewLibraryRepository(db, logger)
	err = repoLibrary.InitWordsMap()
	if err != nil {
		logger.Fatal(err)
	}

	tmpls, err := webtemplate.InitializeTemplates(logger)
	if err != nil {
		logger.Fatal(err)
	}

	hashDB := datastore.InitHashDB()

	r := registry.NewRegistry(db, hashDB, logger, cfg, tmpls, sender)
	e := echo.New()
	e.Validator = validator.NewValidator(logger)

	e = router.NewRouter(e, r.NewAppController(), cfg.Server.SecretKey, tmpls)
	logger.Infof("Server listen at http://%s:%s", cfg.Server.Host, cfg.Server.AppPort)
	if err := e.Start(":" + cfg.Server.AppPort); err != nil {
		logger.Fatalln(err)
	}

}
