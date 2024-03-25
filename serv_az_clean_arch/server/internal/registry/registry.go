package registry

import (
	"server/internal/config"
	"server/internal/infrastructure/webtemplate.go"
	"server/internal/interface/controller"
	"server/internal/interface/repository"
	"server/internal/usercase/interactor"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type registry struct {
	log    *logrus.Logger
	db     *gorm.DB
	config *config.Config
	tmpls  *webtemplate.WebTemplates
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *gorm.DB, log *logrus.Logger, config *config.Config, tmpls *webtemplate.WebTemplates) Registry {
	return &registry{db: db, log: log, config: config, tmpls: tmpls}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.AppController{
		HandlerController: r.NewUserController(),
	}
}

func (r *registry) NewUserController() controller.HandleController {
	userInteractor := interactor.NewUserInteractor(
		repository.NewUserRepository(r.db, r.log),
		repository.NewWordsRepository(r.db, r.log),
	)
	libInteractor := interactor.NewLibraryInteractor(
		repository.NewLibraryRepository(r.db, r.log),
	)

	return controller.NewUserController(userInteractor, libInteractor, r.log, r.config, r.tmpls)
}
