/*
 * @Author: your name
 * @Date: 2021-11-02 19:16:45
 * @LastEditTime: 2021-11-08 17:02:30
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/Router/apiserver.go
 */
package router

import (
	"envelop-rain/common"
	"envelop-rain/configs"
	"envelop-rain/constant"
	"envelop-rain/middleware"
	db "envelop-rain/repository"
	"fmt"
	"time"

	"github.com/bits-and-blooms/bloom"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type APIServer struct {
	sysconfig   configs.SystemConfig
	redisdb     *redis.Client
	mysql       *gorm.DB
	sendall     bool
	bloomFilter *bloom.BloomFilter
	producer    *middleware.KafkaProducer
	consumer    *middleware.KafkaConsumer
}

var server APIServer

func init() {
	// set random seed
	common.SetRandomSeed()
	// get system config
	server.sysconfig = configs.GenerateConfigFromFile()
	// get redis client connection
	server.redisdb = db.GetRedisClient()
	// get mysql connection
	server.mysql = db.GetMySQLCursor()
	server.sendall = false
	// init the bloom filter
	server.bloomFilter = bloom.NewWithEstimates(1000000, 0.01)

	// generate tables
	db.GenerateTables(server.mysql)
	db.SetRemainToDB(server.sysconfig.TotalNum, server.sysconfig.TotalMoney, server.mysql)
	// set config to redis
	configs.SetConfigToRedis(&server.sysconfig, server.redisdb)

	// init kafka producer and consumer
	server.producer = middleware.GetKafkaProducer(configs.GlobalConfig.GetString("Kafka.Topic"))
	server.consumer = middleware.GetKafkaConsumer()
	server.consumer.StartConsume(configs.GlobalConfig.GetString("Kafka.Topic"), server.mysql)
}

func APIServerRun() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	lmForSnatch := middleware.NewRateLimiter(time.Second, configs.GlobalConfig.GetInt64("limiter.SnatchPerSecond"), constant.REQUEST_SNATCH)
	lmForOpen := middleware.NewRateLimiter(time.Minute, configs.GlobalConfig.GetInt64("limiter.OpenPerMinute"), constant.REQUEST_OPEN)
	lmForWallet := middleware.NewRateLimiter(time.Minute, configs.GlobalConfig.GetInt64("limiter.WalletPerMinute"), constant.REQUEST_GETWL)
	r.POST("/snatch", lmForSnatch.Middleware(), SnatchHandler)
	r.POST("/open", lmForOpen.Middleware(), OpenHandler)
	r.POST("/get_wallet_list", lmForWallet.Middleware(), WalletListHandler)
	r.POST("/change_configs", ChangeConfigsHandler)
	r.GET("/flush", FlushDBHandler)
	r.Run()
}

func APIServerStop() {
	fmt.Println("Stop")
}
