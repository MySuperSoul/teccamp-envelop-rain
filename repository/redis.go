/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-07 12:06:39
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/db/redis.go
 */
package repository

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

	// empty the content in redis
	client.FlushDB()
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

func GetRedPacketsByUID(redisdb *redis.Client, uid string) ([]*RedPacket, error) {
	var packets []*RedPacket
	packet_ids, _ := redisdb.LRange(uid+"-wallet", 0, -1).Result()
	for i := len(packet_ids) - 1; i >= 0; i-- {
		packet_id := packet_ids[i]
		vals, _ := redisdb.HGetAll(packet_id).Result()
		packets = append(packets, &RedPacket{
			PacketID:  common.ConvertString(packet_id, "int64").(int64),
			UserID:    common.ConvertString(vals["userid"], "int32").(int32),
			Value:     common.ConvertString(vals["value"], "int32").(int32),
			Opened:    common.ConvertString(vals["opened"], "bool").(bool),
			Timestamp: common.ConvertString(vals["timestamp"], "int64").(int64),
		})
	}

	return packets, nil
}

func GenerateDecrScript() *redis.Script {
	return redis.NewScript(`
	local kc=tonumber(redis.call('GET',KEYS[1]))
	local kmoney=tonumber(redis.call('GET',KEYS[2])) 
	
	if kc==nil or kmoney==nil
	then 
		return -1
	end
	
	local newkc=kc - 1
	local newkmoney=kmoney - ARGV[1]
	
	if newkc < 0 or newkmoney < 0
	then
		return 0
	else
		redis.call('SET',KEYS[1],newkc)
		redis.call('SET',KEYS[2],newkmoney)
		return 1
	end 
	`)
}

func GenerateUserAmountDecrScript() *redis.Script {
	return redis.NewScript(`
	local kUserAmount=tonumber(redis.call('HGET',KEYS[1],KEYS[2]))

	if kUserAmount==nil
	then 
		return -1
	end

	local newkUserAmount=kUserAmount+1

	if newkUserAmount > tonumber(ARGV[1])
	then
		return 0
	else
		redis.call('HSET',KEYS[1],KEYS[2],newkUserAmount)
		return newkUserAmount
	end 
	`)
}
