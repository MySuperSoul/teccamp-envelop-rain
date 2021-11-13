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
	"envelop-rain/middleware"
	db "envelop-rain/repository"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/bits-and-blooms/bloom"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type APIServer struct {
	sysconfig   configs.SystemConfig
	redisdb     *redis.Client
	mysql       *gorm.DB
	sendall     bool
	nomoney     bool
	bloomFilter *bloom.BloomFilter
	producer    *middleware.KafkaProducer
	consumer    *middleware.KafkaConsumer
	logger      *logrus.Logger
	IDGenerator *common.IDProducer
}

var server APIServer

func InitDB() {
	// get redis client connection
	server.redisdb = db.GetRedisClient()
	// get mysql connection
	server.mysql = db.GetMySQLCursor()
	// generate tables
	db.GenerateTables(server.mysql)
	db.SetRemainToDB(&server.sysconfig, server.mysql)
	// set config to redis
	configs.SetConfigToRedis(&server.sysconfig, server.redisdb)
}

func InitLocal() {
	server.sendall = false
	server.nomoney = false
	// init the bloom filter
	server.bloomFilter = bloom.NewWithEstimates(1000000, 0.001)
}

func InitKafka() {
	// init kafka producer and consumer
	server.producer = middleware.GetKafkaProducer(configs.GlobalConfig.GetString("Kafka.Topic"))
	server.consumer = middleware.GetKafkaConsumer()
	server.consumer.StartConsume(configs.GlobalConfig.GetString("Kafka.Topic"), server.mysql)
}

func InitLogger() {
	// init logger
	server.logger = logrus.New()
	server.logger.SetOutput(
		&lumberjack.Logger{
			Filename: configs.GlobalConfig.GetString("Logger.LogPath"),
			Compress: true,
		},
	)
	server.logger.SetFormatter(&logrus.TextFormatter{})
	server.logger.SetLevel(logrus.ErrorLevel)
}

func InitHystrixBreaker() {
	middleware.ConfigHystrix("snatch", configs.GlobalConfig.GetInt("limiter.SnatchPerSecond"))
	middleware.ConfigHystrix("open", configs.GlobalConfig.GetInt("limiter.OpenPerSecond"))
	middleware.ConfigHystrix("wallet", configs.GlobalConfig.GetInt("limiter.WalletPerSecond"))
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
}

func InitIDGenerator() {
	server.IDGenerator = common.NewProducer(0, 0, configs.GlobalConfig.GetInt("common.TotalNum"))
	go server.IDGenerator.StartProducePacketID()
}

func init() {
	// set random seed
	common.SetRandomSeed()
	// get system config
	server.sysconfig = configs.GenerateConfigFromFile()

	InitDB()
	InitLocal()
	InitKafka()
	InitLogger()
	InitHystrixBreaker()
	InitIDGenerator()
}

func APIServerRun() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	BreakerSnatch := &middleware.HystrixMiddleWare{Name: "snatch"}
	BreakerOpen := &middleware.HystrixMiddleWare{Name: "open"}
	BreakerWallet := &middleware.HystrixMiddleWare{Name: "wallet"}

	r.POST("/snatch", BreakerSnatch.Middleware(), SnatchHandler)
	r.POST("/open", BreakerOpen.Middleware(), OpenHandler)
	r.POST("/get_wallet_list", BreakerWallet.Middleware(), WalletListHandler)
	r.POST("/change_configs", ChangeConfigsHandler)
	r.Run()
}

func APIServerStop() {
	fmt.Println("Stop")
}
