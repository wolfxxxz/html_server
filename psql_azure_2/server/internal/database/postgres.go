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

//https://gorm.io/ru_RU/docs/connecting_to_the_database.html

//sqlDB.SetConnMaxLifetime(time.Minute * 5)
//db.DB().SetMaxIdleConns(10)
//db.DB().SetMaxOpenConns(100)

/*
// var db *sql.DB
var server = "webservertarget.database.windows.net"
var port = 1433
var user = "whitecat"
var password = "BlackCat2"
var database = "psql_web_server"
*/

// old local connection

/*
	sqlDB, err := db.DB()
	if err != nil {
		appErr := apperrors.SetupDatabaseErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}

	tNum, err := strconv.Atoi(conf.Postgres.TimeoutQuery)
	if err != nil {
		appErr := apperrors.SetupDatabaseErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}

	dsnWithoutPassword := fmt.Sprintf("%v://%v:%v/%v?&user=%v&password=[great secret]&dbname=%v&TimeZone=%s",
		conf.Postgres.SqlType, conf.Postgres.SqlHost, conf.Postgres.SqlPort, conf.Postgres.SqlType,
		conf.Postgres.UserName, conf.Postgres.DBName, conf.Postgres.TimeZone,
	)
	log.Infof("Trying to connect to Postgres.\n %s", dsnWithoutPassword)
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(tNum))
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		appErr := apperrors.SetupDatabaseErr.AppendMessage(fmt.Sprintf("PingErr %v", err))
		log.Error(appErr)
		return nil, appErr
	}
*/

// dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
// 	conf.Postgres.SqlHost, conf.Postgres.UserName, conf.Postgres.Password, conf.Postgres.SqlPort, conf.Postgres.DBName)
