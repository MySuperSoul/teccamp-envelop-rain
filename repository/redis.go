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

func GetRedPacketsByUID(redisdb *redis.Client, uid string) ([]*RedPacket, int32) {
	var packets []*RedPacket
	balance := int32(0)
	packet_ids, _ := redisdb.LRange(uid+"-wallet", 0, -1).Result()
	for i := len(packet_ids) - 1; i >= 0; i-- {
		packet_id := packet_ids[i]
		vals, _ := redisdb.HGetAll("packet-" + packet_id).Result()
		packet := RedPacket{
			PacketID:  common.ConvertString(packet_id, "int64").(int64),
			UserID:    common.ConvertString(vals["userid"], "int32").(int32),
			Value:     common.ConvertString(vals["value"], "int32").(int32),
			Opened:    common.ConvertString(vals["opened"], "bool").(bool),
			Timestamp: common.ConvertString(vals["timestamp"], "int64").(int64),
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
// ARGV: [uidstr, timestamp, total_num, max_amount, p]
// return: 0: success, -1: not lucky, -2: no packet, -3: user exceed (And cur num)

func SnatchScript() *redis.Script {
	return redis.NewScript(`
		local timestamp = tonumber(ARGV[2])
		local total_num = tonumber(ARGV[3])
		local current_num = tonumber(redis.call('GET', KEYS[1]))

		if current_num == total_num
		then 
			return -2
		end

		local wallet_name = ARGV[1] .. "-wallet"
		local useramount = tonumber(redis.call('LLEN', wallet_name))
		local maxamount = tonumber(ARGV[4])
		if useramount == maxamount 
		then
			return -3
		end

		local p = tonumber(ARGV[5])
		math.randomseed(timestamp)
		local up = math.random()
		if up > p
		then
			return -1
		end

		redis.call('SET', KEYS[1], current_num + 1)
		redis.call('LPUSH', wallet_name, timestamp)
		return current_num + 1
	`)
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

// KEYS: ["RemainNum", "RemainMoney"]
// ARGV: [MaxMoney, MinMoney, TimeStamp]
func GeneratePacketScript() *redis.Script {
	return redis.NewScript(`
	math.randomseed(ARGV[3])
	local remain_num = tonumber(redis.call('GET', KEYS[1]))
	local remain_money = tonumber(redis.call('GET',KEYS[2]))
	local value = 0

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

func GenerateOpenPacketScript() *redis.Script {
	return redis.NewScript(`
	local function toboolean(str)
    	local bool = false
    	if str == "true" then
        	bool = true
    	end
    	return bool
	end
	local kOpened=toboolean(redis.call('HGET',KEYS[1],KEYS[2]))
	if kOpened==nil
	then 
		return -1
	end

	if kOpened==true
	then
		return 0
	else
		local kValue=tonumber(redis.call('HGET',KEYS[1],KEYS[3]))
		redis.call('HINCRBY',KEYS[4],KEYS[5],kValue)
		redis.call('HSET',KEYS[1],KEYS[2],'true')
		return kValue
	end
	`)
}

func GenerateChangeScript() *redis.Script {
	return redis.NewScript(`
	local kc=tonumber(redis.call('GET',KEYS[1]))
	local kmoney=tonumber(redis.call('GET',KEYS[2])) 

	if kc==nil or kmoney==nil
	then 
		return -1
	end

	local newkc=kc-ARGV[1]
	local newkmoney=kmoney - ARGV[2]

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
