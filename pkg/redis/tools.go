package redis

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var (
	client *redis.Client
	mutex  = &sync.Mutex{}
)

func LoadRedis(pass, host string, port, db int) error {
	mutex.Lock()
	defer mutex.Unlock()

	add := fmt.Sprintf("%s:%d", host, port)
	if len(pass) == 0 {
		logrus.WithFields(logrus.Fields{
			"add": add,
		}).Debug("Redis connection without password")
		client = redis.NewClient(&redis.Options{
			Addr: add,
			DB:   db,
		})
	} else {
		logrus.WithFields(logrus.Fields{
			"add": add,
		}).Debug("Redis connection with password")
		client = redis.NewClient(&redis.Options{
			Addr:     add,
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
