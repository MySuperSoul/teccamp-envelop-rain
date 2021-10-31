package main

import (
	"envelop-rain/common"
	"envelop-rain/configs"
	"envelop-rain/db"
	"envelop-rain/redpacket"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var sysconfig configs.SystemConfig
var redisdb *redis.Client
var mysql *gorm.DB

func InitServer() {
	// set random seed
	common.SetRandomSeed()
	// get system config
	sysconfig = configs.GenerateConfigFromFile(CONFIG_PATH)
	// get redis client connection
	redisdb = db.GetRedisClient()
	// get mysql connection
	mysql = db.GetMySQLCursor()
	// generate tables
	db.GenerateTables(mysql)
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
	log.Info("snatch by user: %s", uid)

	// first to judge whether has packet left
	remain_num := configs.GetSingleValueFromRedis(redisdb, "RemainNum", "int32").(int32)
	remain_money := configs.GetSingleValueFromRedis(redisdb, "RemainMoney", "float32").(float32)

	if remain_num == 0 {
		c.JSON(200, gin.H{"code": SNATCH_NO_RED_PACKET, "msg": SNATCH_NO_RED_PACKET_MESSAGE, "data": gin.H{}})
		return
	}

	// Then to judge whether the user is lucky enough
	if p := configs.GetSingleValueFromRedis(redisdb, "P", "float32").(float32); common.Rand() > p {
		c.JSON(200, gin.H{"code": SNATCH_NOT_LUCKY, "msg": SNATCH_NOT_LUCKY_MESSAGE, "data": gin.H{}})
		return
	}

	// Then perform later operations
	// First generate the red packet
	packet := redpacket.GetRedPacket(remain_num, remain_money, sysconfig.MinMoney, sysconfig.MaxMoney)
	packet.UserID = common.ConvertString(uid, "int32").(int32)

	// get user information and update to database
	user := db.User{UserID: packet.UserID, Amount: 0, Balance: 0.}
	mysql.Where(db.User{UserID: packet.UserID}).FirstOrCreate(&user)

	// update remain value to redis
	redisdb.Decr("RemainNum")
	redisdb.Set("RemainMoney", remain_money-packet.Value, 0)

	// send message
	c.JSON(200, gin.H{
		"code": SNATCH_SUCCESS,
		"msg":  SNATCH_SUCCESS_MESSAGE,
		"data": gin.H{},
	})

	// insert into database
}

func OpenHandler(c *gin.Context) {

}

func WalletListHandler(c *gin.Context) {

}
