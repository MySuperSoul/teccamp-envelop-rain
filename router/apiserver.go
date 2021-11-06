/*
 * @Author: your name
 * @Date: 2021-11-02 19:16:45
 * @LastEditTime: 2021-11-06 22:44:15
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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type APIServer struct {
	sysconfig configs.SystemConfig
	redisdb   *redis.Client
	mysql     *gorm.DB
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
	// generate tables
	db.GenerateTables(server.mysql)
	// set config to redis
	configs.SetConfigToRedis(&server.sysconfig, server.redisdb)
}

func APIServerRun() {
	r := gin.Default()
	lmForSnatch := middleware.NewRateLimiter(time.Second, configs.GlobalConfig.GetInt64("limiter.SnatchPerSecond"))
	lmForOpen := middleware.NewRateLimiter(time.Minute, configs.GlobalConfig.GetInt64("limiter.OpenPerMinute"))
	lmForWallet := middleware.NewRateLimiter(time.Minute, configs.GlobalConfig.GetInt64("limiter.WalletPerMinute"))
	r.POST("/snatch", lmForSnatch.Middleware(), SnatchHandler)
	r.POST("/open", lmForOpen.Middleware(), OpenHandler)
	r.POST("/get_wallet_list", lmForWallet.Middleware(), WalletListHandler)
	r.Run()
}

func APIServerStop() {
	fmt.Println("Stop")
}
