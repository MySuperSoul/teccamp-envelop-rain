/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-07 17:15:14
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/db/redis.go
 */
package repository

import (
	"envelop-rain/common"
	"envelop-rain/configs"
	"os"
	"sort"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func GetRedPacketsByUID(redisdb *redis.Client, mysql *gorm.DB, uid string) ([]*RedPacket, int32) {
	var packets []*RedPacket
	balance := int32(0)
	packet_ids, _ := redisdb.LRange(uid+"-wallet", 0, -1).Result()
	for i := len(packet_ids) - 1; i >= 0; i-- {
		var packet RedPacket
		packet_id := packet_ids[i]
		if n, _ := redisdb.Exists("packet-" + packet_id).Result(); n == 0 { // not in redis
			// Then fetch db and set to redis
			mysql.Where("packet_id = ?", common.ConvertString(packet_id, "int64").(int64)).Find(&packet)
			redisdb.HMSet("packet-"+packet_id, map[string]interface{}{
				"userid":    packet.UserID,
				"value":     packet.Value,
				"opened":    packet.Opened,
				"timestamp": packet.Timestamp,
			})
		} else { // in redis cache
			vals, _ := redisdb.HGetAll("packet-" + packet_id).Result()
			packet = RedPacket{
				PacketID:  common.ConvertString(packet_id, "int64").(int64),
				UserID:    common.ConvertString(vals["userid"], "int32").(int32),
				Value:     common.ConvertString(vals["value"], "int32").(int32),
				Opened:    common.ConvertString(vals["opened"], "bool").(bool),
				Timestamp: common.ConvertString(vals["timestamp"], "int64").(int64),
			}
		}
		packets = append(packets, &packet)
		if packet.Opened {
			balance += int32(packet.Value)
		}
	}

	sort.SliceStable(packets, func(i, j int) bool {
		return packets[i].Timestamp < packets[j].Timestamp
	})

	return packets, balance
}

// keys: ["CurrentNum", ]
// ARGV: [uidstr, timestamp, total_num, max_amount]
// return: 0: success, -1: not lucky, -2: no packet, -3: user exceed (And cur num)

func SnatchScript() *redis.Script {
	return redis.NewScript(`
		local timestamp = tonumber(ARGV[2])
		local total_num = tonumber(ARGV[3])
		local current_num = tonumber(redis.call('GET', KEYS[1]))

		if current_num >= total_num
		then 
			return -2
		end

		local wallet_name = ARGV[1] .. "-wallet"
		local useramount = tonumber(redis.call('LLEN', wallet_name))
		local maxamount = tonumber(ARGV[4])
		if useramount >= maxamount 
		then
			return -3
		end

		redis.call('SET', KEYS[1], current_num + 1)
		redis.call('LPUSH', wallet_name, timestamp)
		return useramount + 1
	`)
}

// KEYS: ["RemainNum", "RemainMoney"]
// ARGV: [MaxMoney, MinMoney, TimeStamp]
func GeneratePacketScript() *redis.Script {
	return redis.NewScript(`
	math.randomseed(ARGV[3])
	local remain_num = tonumber(redis.call('GET', KEYS[1]))
	local remain_money = tonumber(redis.call('GET',KEYS[2]))
	local value = 0

	if remain_num < 0 or remain_money < 0
	then
		return -1
	end

	if remain_num == 1
	then
		value = math.min(ARGV[1], remain_money)
	else
		local mean_money = math.floor(remain_money / remain_num)
		local max_money = math.min(ARGV[1], 2*mean_money-ARGV[2])
		value = ARGV[2] + math.random(max_money - ARGV[2] + 1) - 1
	end

	redis.call('SET', KEYS[1], remain_num - 1)
	redis.call('SET', KEYS[2], remain_money - value)
	return value
	`)
}

func GenerateChangeScript() *redis.Script {
	return redis.NewScript(`
	local kmoney=tonumber(redis.call('GET',KEYS[1])) 
	if kmoney==nil
	then 
		return -1
	end
	local newkmoney=kmoney - ARGV[1]
	redis.call('SET',KEYS[1],newkmoney)
	return 1
	`)
}
