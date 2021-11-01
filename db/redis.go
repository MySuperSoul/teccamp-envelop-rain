/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-01 16:35:34
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/db/redis.go
 */
package db

import (
	"envelop-rain/common"
	"envelop-rain/configs"
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// redis config
type RedisConfig struct {
	RedisAddr, RedisPassword string
	RedisPoolSize            int
}

func initRedisConfig() RedisConfig {
	var redis_config RedisConfig = RedisConfig{
		configs.GlobalConfig.GetString("Redis.RedisAddr"),
		configs.GlobalConfig.GetString("Redis.RedisPassword"),
		configs.GlobalConfig.GetInt("Redis.RedisPoolSize")}
	return redis_config
}

func GetRedisClient() *redis.Client {
	redis_config := initRedisConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     redis_config.RedisAddr,
		Password: redis_config.RedisPassword, // no password set
		DB:       0,                          // use default DB
		PoolSize: redis_config.RedisPoolSize,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("Redis connection failed")
		os.Exit(0)
	}
	log.Info("Redis connection success")
	return client
}

func GetSingleValueFromRedis(redisdb *redis.Client, key string, datatype string) interface{} {
	val, err := redisdb.Get(key).Result()
	if err != nil {
		log.Fatalf("Read from key: %s failed", key)
		panic(err)
	}

	return common.ConvertString(val, datatype)
}
