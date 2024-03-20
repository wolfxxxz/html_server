package database

import (
	"context"
	"fmt"
	"server/internal/apperrors"
	"server/internal/config"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB interface {
	SetupDatabase(ctx context.Context, conf *config.Config, logger *logrus.Logger) (*gorm.DB, error)
	PingEveryMinuts(ctx context.Context, timeOutMinute int, log *logrus.Logger) error
}

type postgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB() PostgresDB {
	return &postgresDB{}
}

func (p *postgresDB) SetupDatabase(ctx context.Context, conf *config.Config, log *logrus.Logger) (*gorm.DB, error) {
	if conf.Postgres.DBName == "" {
		appErr := apperrors.SetupDatabaseErr.AppendMessage("config DBName is empty")
		log.Error(appErr)
		return nil, appErr
	}

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		conf.Postgres.UserName, conf.Postgres.Password, conf.Postgres.SqlHost, conf.Postgres.SqlPort, conf.Postgres.DBName)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //LogMode(logger.Silent)
	})

	if err != nil {
		appErr := apperrors.SetupDatabaseErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}

	p.DB = db

	//keep alive
	sqlDB, err := db.DB()
	if err != nil {
		appErr := apperrors.SetupDatabaseErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	err = sqlDB.Ping()
	if err != nil {
		appErr := apperrors.SetupDatabaseErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}
	//

	log.Info("DB Postgres has been connected, DB.Ping success ")
	return db, nil
}

func (p *postgresDB) PingEveryMinuts(ctx context.Context, timeOutMinute int, log *logrus.Logger) error {
	ticker := time.NewTicker(time.Minute * time.Duration(timeOutMinute))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			appErr := apperrors.PingEveryMinutsErr.AppendMessage(ctx.Err())
			log.Info(appErr)
			return appErr
		case <-ticker.C:
			sqlDB, err := p.DB.DB()
			if err != nil {
				appErr := apperrors.PingEveryMinutsErr.AppendMessage(err)
				log.Error(appErr)
			}

			err = sqlDB.Ping()
			if err != nil {
				appErr := apperrors.PingEveryMinutsErr.AppendMessage(err)
				log.Error(appErr)
			}

			log.Info("Ping ring rang, DB woke up")
		}
	}

}
