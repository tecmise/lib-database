package database

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	tools "github.com/tecmise/lib-database/pkg/redis"
)

type RedisRepository struct {
	dbRedis *redis.Client
	once    sync.Once
}

type RedisConfiguration struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName int
}

var Redis = &RedisRepository{}

func (r *RedisRepository) Start(configuration RedisConfiguration) {
	_ = tools.LoadRedis(
		configuration.DBPass,
		configuration.DBHost,
		configuration.DBPort,
		configuration.DBName,
	)
}

func (r *RedisRepository) Stop() {
	if r.dbRedis != nil {
		_ = r.dbRedis.Close()
	}
}

func (r *RedisRepository) GetInstance() *redis.Client {
	r.once.Do(func() {
		var err error
		r.dbRedis, err = tools.GetRedis()
		if err != nil {
			panic(err.Error())
		}
	})
	return r.dbRedis
}

func (r *RedisRepository) Ping() error {
	ctx := context.Background()
	client := r.GetInstance()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
