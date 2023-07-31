package cache

import (
	"context"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var doConnectionOnce sync.Once

var cacheClient *redis.Client

func getConnection() *redis.Client {
	doConnectionOnce.Do(func() {
		opt, err := redis.ParseURL(commons.GoDotEnvVariable("CACHE_CONNECTION"))
		if err != nil {
			panic(err)
		}
		cacheClient = redis.NewClient(opt)
	})

	return cacheClient
}

func Get(key string) (string, error) {
	client := getConnection()
	return client.Get(context.Background(), key).Result()
}

func Set(key string, value string, expires time.Duration) {
	client := getConnection()
	err := client.Set(context.Background(), key, value, expires).Err()
	if err != nil {
		panic(err)
	}
}

func HSet(key string, value map[string]string) {
	client := getConnection()
	for k, val := range value {
		err := client.HSet(context.Background(), key, k, val)
		if err != nil {
			panic(err)
		}
	}
}

func HGet(key string) map[string]string {
	client := getConnection()
	collection := client.HGetAll(context.Background(), key)
	return collection.Val()
}

func Remember(key string, cmd func() (string, error), expires time.Duration) (string, error) {
	val, err := Get(key)
	if err == nil && val != "" {
		return val, nil
	}

	val, err = cmd()
	if err != nil {
		return "", err
	}
	Set(key, val, expires)

	return val, err
}

func RememberForever(key string, cmd func() (string, error)) (string, error) {
	return Remember(key, cmd, 0)
}
