package database

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	tools "github.com/tecmise/lib-database/pkg/gorm"
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

// MySQL var to use
var Postgres = &PostgresRepository{}

// StartPostgres start the DB
func (r *PostgresRepository) Start(configuration PostgresConfiguration) {
	_ = tools.LoadGormPostGres(
		configuration.DBUser,
		configuration.DBPass,
		configuration.DBHost,
		configuration.DBPort,
		configuration.DBName,
		configuration.Ssl)
}

// StopPostgres stop the DB
func (r *PostgresRepository) Stop() {
	defer r.dbPostgres.Close()
}

// GetInstance returns a unique instance of gorm.DB
func (r *PostgresRepository) GetInstance() *gorm.DB {

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
		r.dbPostgres.SingularTable(true)
		r.dbPostgres.LogMode(logLevel)
	})
	return r.dbPostgres
}
