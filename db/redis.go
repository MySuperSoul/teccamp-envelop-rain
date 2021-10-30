package db

import (
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// redis config
const (
	REDIS_ADDR      string = "localhost:6379"
	REDIS_PASSWORD  string = ""
	REDIS_POOL_SIZE int    = 100
)

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR,
		Password: REDIS_PASSWORD, // no password set
		DB:       0,              // use default DB
		PoolSize: REDIS_POOL_SIZE,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("Redis connection failed")
		os.Exit(0)
	}
	log.Info("Redis connection success")
	return client
}
