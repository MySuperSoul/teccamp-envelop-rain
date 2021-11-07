/*
 * @Author: your name
 * @Date: 2021-11-02 19:16:45
 * @LastEditTime: 2021-11-02 21:12:29
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/Router/apiserver.go
 */
package router

import (
	"envelop-rain/common"
	"envelop-rain/configs"
	db "envelop-rain/repository"
	"fmt"

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
	r.POST("/snatch", SnatchHandler)
	r.POST("/open", OpenHandler)
	r.POST("/get_wallet_list", WalletListHandler)
	r.GET("/flush", FlushDBHandler)
	r.Run()
}

func APIServerStop() {
	fmt.Println("Stop")
}
