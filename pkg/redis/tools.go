package redis

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
	mutex  = &sync.Mutex{}
)

func LoadRedis(user, pass, host string, port, db int) error {
	mutex.Lock()
	defer mutex.Unlock()

	if len(pass) == 0 {
		client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", host, port),
			DB:   db,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: pass,
			DB:       db,
		})
	}

	return nil
}

func GetRedis() (*redis.Client, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if client != nil {
		return client, nil
	}

	return nil, fmt.Errorf("redis client not found")
}
