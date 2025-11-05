package tools

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tecmise/lib-database/pkg/logger"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"net/url"
	"os"
)

type singleton struct {
	Connection *gorm.DB
}

// MultipleDatabaseOnPoolError godoc
type MultipleDatabaseOnPoolError struct{}

// PoolWithoutInstanceError godoc
type PoolWithoutInstanceError struct{}

// Vars Gorm godoc
var (
	poolGormDb    map[string]*gorm.DB = make(map[string]*gorm.DB)
	LogModeEnable                     = func() bool {
		return os.Getenv("LOG_MODE") == "true"
	}()
)

// LoadGormPostgres with the following parameters:
/**
 * user string, pass string, host string, port int, dbName string
**/
func LoadGormPostgres(user string, pass string, host string, port int, dbName string, sslMode bool, schema string) error {
	return LoadGorm("postgres", user, pass, host, port, dbName, sslMode, schema, postgres.Open)
}

// LoadGormPostgres with the following parameters:
/**
 * user string, pass string, host string, port int, dbName string
**/
func LoadGormClickhouse(user string, pass string, host string, port int, dbName string, sslMode bool) error {
	return LoadGorm("clickhouse", user, pass, host, port, dbName, sslMode, "", clickhouse.Open)
}

// LoadGorm with the following parameters:
/**
 * driverName string (such as mysql) user string, pass string, host string, port int, dbName string
**/
func LoadGorm(driverName string, user string, pass string, host string, port int, dbName string, sslMode bool, schema string, dialector func(dsn string) gorm.Dialector) error {
	var err error

	if poolGormDb[dbName] == nil {
		poolGormDb[dbName], err = getGormConnection(driverName, user, pass, host, port, dbName, sslMode, schema, dialector)
	}

	return err
}

// GetGormDb return gorm.DB instance
func GetGormDb(dbNameParam ...string) (*gorm.DB, error) {
	dbName, err := defineDatabaseName(dbNameParam, len(poolGormDb), func() string {
		return firstKeyFromGormPool(poolGormDb)
	})

	if err != nil {
		return nil, err
	}
	if poolGormDb[dbName] == nil {
		return nil, errors.New("LoadGorm/SetGormDb wasn't called for database: " + dbName)
	}

	return poolGormDb[dbName], nil
}

// SetGormDb godoc
func SetGormDb(gormDb *gorm.DB, dbName string) {
	poolGormDb[dbName] = gormDb
}

func getGormConnection(driverName string, user string, pass string, host string, port int, dbName string, sslMode bool, schema string, dialector func(dsn string) gorm.Dialector) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	var dsn string
	dsn, err = generateDsn(driverName, user, pass, host, port, dbName, sslMode, schema)

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	db, err = gorm.Open(dialector(dsn), &gorm.Config{
		Logger: logger.NewGormLogrus(log, gormLogger.Info),
	})

	if err != nil {
		logrus.WithFields(map[string]interface{}{
			"host":    host,
			"port":    port,
			"db_name": dbName,
			"ssl":     sslMode,
		}).Error("Erro ao subir conexao")
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	//db.DB().SetConnMaxLifetime(time.Duration(getenv.DbConnPoolLifeTime) * time.Minute)
	//db.DB().SetMaxIdleConns(getenv.DbConnPoolMaxIdle)
	//db.DB().SetMaxOpenConns(getenv.DbConnPoolMaxOpen)
	return db, err
}

func firstKeyFromGormPool(object map[string]*gorm.DB) string {
	for k := range object {
		return k
	}

	return ""
}

// Error godoc
func (e *MultipleDatabaseOnPoolError) Error() string {
	return fmt.Sprintf("It isn't allowed define a default database. You should pass the database name instead.")
}

// Error godoc
func (e *PoolWithoutInstanceError) Error() string {
	return fmt.Sprintf("Can't define a default database. You should Set or Load a instance first.")
}

func defineDatabaseName(dbNameParam []string, poolSize int, firstKeyOfPool func() string) (string, error) {
	var dbName string

	if len(dbNameParam) == 0 {
		if poolSize == 0 {
			return "", &PoolWithoutInstanceError{}
		}

		if poolSize > 1 {
			return "", &MultipleDatabaseOnPoolError{}
		}

		dbName = firstKeyOfPool()
	} else {
		dbName = dbNameParam[0]
	}

	return dbName, nil
}

func generateDsn(driverName string, user string, pass string, host string, port int, dbName string, sslMode bool, schema string) (string, error) {
	var sslOption string

	switch driverName {
	case "clickhouse":
		// Escapa user/pass/dbName para evitar caracteres problemáticos na URL
		return fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s",
			url.PathEscape(user),
			url.PathEscape(pass),
			host,
			port,
			url.PathEscape(dbName),
		), nil
	case "postgres", "postgresql":
		if sslMode {
			sslOption = "verify-ca"
		} else {
			sslOption = "disable"
		}

		if schema != "" {
			return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s search_path=%s",
				host, port, user, dbName, pass, sslOption, schema,
			), nil
		}

		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			host, port, user, dbName, pass, sslOption,
		), nil
	default:
		return "", errors.New("unsupported driver: " + driverName)
	}

}
