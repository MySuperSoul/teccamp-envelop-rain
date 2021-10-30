package main

import (
	"envelop-rain/common"
	"envelop-rain/configs"
	"envelop-rain/db"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var sysconfig configs.SystemConfig
var redisdb *redis.Client

func InitServer() {
	// set random seed
	common.SetRandomSeed()
	// get system config
	sysconfig = configs.GenerateConfigFromFile(CONFIG_PATH)
	// get redis client connection
	redisdb = db.GetRedisClient()
	// set config to redis
	configs.SetConfigToRedis(&sysconfig, redisdb)
}

func main() {
	// init server
	InitServer()

	// start server
	r := gin.Default()
	r.POST("/snatch", SnatchHandler)
	r.POST("/open", OpenHandler)
	r.POST("/get_wallet_list", WalletListHandler)
	r.Run()
}

func SnatchHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Info("snatch by user: %d", uid)

	// first to judge whether has packet left

	// Then to judge whether the user is lucky enough

	// Then perform later operations, user get this packet
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{},
	})

}

func OpenHandler(c *gin.Context) {

}

func WalletListHandler(c *gin.Context) {

}
