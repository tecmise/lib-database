package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tecmise/lib-database/pkg/database"
)

// BaseRepository base for all repositories
type BaseRepository struct{}

// conn return a database connection
func (b BaseRepository) Conn() *gorm.DB {
	return database.Postgres.GetInstance()
}

func (b BaseRepository) NextVal(sequence string) int64 {
	var next int64
	command := fmt.Sprintf("select nextval('%s')", sequence)
	b.Conn().Raw(command).Row().Scan(&next)
	return next
}
