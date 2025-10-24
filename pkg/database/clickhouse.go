package database

import (
	"github.com/sirupsen/logrus"
	tools "github.com/tecmise/lib-database/pkg/gorm"
	"gorm.io/gorm"
	"os"
	"strings"
	"sync"
)

type ClickhouseRepository struct {
	dbPostgres *gorm.DB
	once       sync.Once
}

type ClickhouseConfiguration struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string
	Ssl    bool
}

var Clickhouse = &ClickhouseRepository{}

func (r *ClickhouseRepository) Start(configuration ClickhouseConfiguration) {
	_ = tools.LoadGormClickhouse(
		configuration.DBUser,
		configuration.DBPass,
		configuration.DBHost,
		configuration.DBPort,
		configuration.DBName,
		configuration.Ssl)
}

func (r *ClickhouseRepository) Stop() {
	if r.dbPostgres != nil {
		sqlDB, err := r.dbPostgres.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

func (r *ClickhouseRepository) GetInstance() *gorm.DB {

	r.once.Do(func() {
		logLevel := os.Getenv("SHOW_SQL") != "" && strings.ToLower(os.Getenv("SHOW_SQL")) == "true"
		if !logLevel {
			logrus.Debugf("Log Level Disabled")
			logrus.Debugf("To Enable sql mode set variable SHOW_SQL to true")
		}
		var err error
		r.dbPostgres, err = tools.GetGormDb()
		if err != nil {
			panic(err.Error())
		}
		if logLevel {
			r.dbPostgres.Logger.LogMode(4)
		}
	})
	return r.dbPostgres
}
