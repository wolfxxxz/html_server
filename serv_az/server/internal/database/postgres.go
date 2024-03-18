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

	//log.Error(dsn)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //LogMode(logger.Silent)
	})
	if err != nil {
		appErr := apperrors.SetupDatabaseErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}

	//keep alive
	sqlDB, err := db.DB()
	if err != nil {
		appErr := apperrors.SetupDatabaseErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	//

	log.Info("DB Postgres has been connected, DB.Ping success ")
	return db, nil
}
