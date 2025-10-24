package database

import (
	"github.com/sirupsen/logrus"
	tools "github.com/tecmise/lib-database/pkg/gorm"
	"gorm.io/gorm"
	"os"
	"strings"
	"sync"
)

type PostgresRepository struct {
	dbPostgres *gorm.DB
	once       sync.Once
}

type PostgresConfiguration struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string
	Ssl    bool
}

var Postgres = &PostgresRepository{}

func (r *PostgresRepository) Start(configuration PostgresConfiguration) {
	_ = tools.LoadGormPostgres(
		configuration.DBUser,
		configuration.DBPass,
		configuration.DBHost,
		configuration.DBPort,
		configuration.DBName,
		configuration.Ssl)
}

func (r *PostgresRepository) Stop() {
	if r.dbPostgres != nil {
		sqlDB, err := r.dbPostgres.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

func (r *PostgresRepository) GetInstance(key string) *gorm.DB {

	r.once.Do(func() {
		logLevel := os.Getenv("SHOW_SQL") != "" && strings.ToLower(os.Getenv("SHOW_SQL")) == "true"
		if !logLevel {
			logrus.Debugf("Log Level Disabled")
			logrus.Debugf("To Enable sql mode set variable SHOW_SQL to true")
		}
		var err error
		r.dbPostgres, err = tools.GetGormDb(key)
		if err != nil {
			panic(err.Error())
		}
		if logLevel {
			r.dbPostgres.Logger.LogMode(4)
		}
	})
	return r.dbPostgres
}
