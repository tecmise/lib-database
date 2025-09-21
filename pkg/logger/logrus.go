package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
)

// GormLogrus implementa gorm/logger.Interface usando Logrus.
type GormLogrus struct {
	logger *logrus.Logger
	level  gormLogger.LogLevel
}

// NewGormLogrus cria um logger para GORM a partir de um logrus.Logger já configurado.
func NewGormLogrus(baseLogger *logrus.Logger, lvl gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogrus{
		logger: baseLogger,
		level:  lvl,
	}
}

// LogMode ajusta o nível de log (interna do GORM)
func (l *GormLogrus) LogMode(lvl gormLogger.LogLevel) gormLogger.Interface {
	new := *l
	new.level = lvl
	return &new
}

func (l *GormLogrus) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormLogger.Info {
		l.logger.Infof(msg, data...)
	}
}

func (l *GormLogrus) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormLogger.Warn {
		l.logger.Warnf(msg, data...)
	}
}

func (l *GormLogrus) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormLogger.Error {
		l.logger.Errorf(msg, data...)
	}
}

func (l *GormLogrus) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	if l.level <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	entry := l.logger.WithFields(logrus.Fields{
		"elapsed": fmt.Sprintf("%v", elapsed),
		"rows":    rows,
	})

	switch {
	case err != nil && l.level >= gormLogger.Error:
		entry.Errorf("%s | %s", err, sql)
	case elapsed > 200*time.Millisecond && l.level >= gormLogger.Warn:
		entry.Warnf("SLOW SQL >200ms: %s", sql)
	case l.level >= gormLogger.Info:
		entry.Infof(sql)
	}
}
