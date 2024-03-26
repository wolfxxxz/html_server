package registry

import (
	"server/internal/config"
	"server/internal/infrastructure/datastore"
	"server/internal/infrastructure/email"
	"server/internal/infrastructure/webtemplate.go"
	"server/internal/interface/controller"
	"server/internal/interface/repository"
	"server/internal/usercase/comparer"
	"server/internal/usercase/interactor"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type registry struct {
	log    *logrus.Logger
	db     *gorm.DB
	hashDB *datastore.HashDB
	config *config.Config
	tmpls  *webtemplate.WebTemplates
	sender email.Sender
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *gorm.DB, hashDB *datastore.HashDB, log *logrus.Logger, config *config.Config, tmpls *webtemplate.WebTemplates, sender email.Sender) Registry {
	return &registry{db: db, hashDB: hashDB, log: log, config: config, tmpls: tmpls, sender: sender}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.AppController{
		HandlerController: r.NewHandlersController(),
	}
}

func (r *registry) NewHandlersController() controller.HandleController {
	userInteractor := interactor.NewUserInteractor(
		repository.NewUserRepository(r.db, r.log),
		repository.NewWordsRepository(r.db, r.log),
		r.sender,
	)
	libInteractor := interactor.NewLibraryInteractor(
		repository.NewLibraryRepository(r.db, r.log),
		repository.NewWordsRepository(r.db, r.log),
		repository.NewBackUpCopyRepo(backupXLS, r.log),
	)
	comparr := comparer.NewComparer(libInteractor, userInteractor, r.log)

	return controller.NewHandlersController(comparr, userInteractor, libInteractor, r.hashDB, r.log, r.config, r.tmpls)
}

const backupXLS = "save_copy/library.xlsx"
