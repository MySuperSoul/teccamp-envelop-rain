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
	uid := common.ConvertString(c.PostForm("uid"), "int32").(int32)
	log.Infof("snatch by user: %d", uid)

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
	// get user information
	user := db.User{UserID: uid, Amount: 0, Balance: 0.}
	mysql.Where(db.User{UserID: uid}).FirstOrCreate(&user)

	// Then to check the maxamount
	max_amount := configs.GetSingleValueFromRedis(redisdb, "MaxAmount", "int32").(int32)
	if user.Amount == max_amount {
		c.JSON(200, gin.H{"code": SNATCH_EXCEED_MAX_AMOUNT, "msg": SNATCH_EXCEED_MAX_AMOUNT_MESSAGE, "data": gin.H{}})
		return
	}

	// First generate the red packet
	packet := redpacket.GetRedPacket(remain_num, remain_money, sysconfig.MinMoney, sysconfig.MaxMoney)
	packet.UserID = uid

	// update remain value to redis
	redisdb.Decr("RemainNum")
	redisdb.Set("RemainMoney", remain_money-packet.Value, 0)

	// update user amount and insert the red packet
	user.Amount++
	mysql.Model(&user).Update("amount", user.Amount)
	mysql.Create(&packet)

	// send message
	c.JSON(200, gin.H{
		"code": SNATCH_SUCCESS,
		"msg":  SNATCH_SUCCESS_MESSAGE,
		"data": gin.H{"envelop_id": packet.PacketID, "max_count": max_amount, "cur_count": user.Amount},
	})
}

func OpenHandler(c *gin.Context) {
	userid := common.ConvertString(c.PostForm("uid"), "int32").(int32)
	packetid := common.ConvertString(c.PostForm("envelop_id"), "int64").(int64)
	log.Infof("Envelop %d opened by %d.", packetid, userid)

	var user db.User
	var packet db.RedPacket
	result := mysql.First(&user, userid)
	if result.RowsAffected == 0 {
		log.Errorf("Invalid user id: %d, block him.", userid)
		c.JSON(200, gin.H{"code": OPEN_INVALID_USER, "msg": OPEN_INVALID_USER_MESSAGE, "data": gin.H{}})
		return
	}
	result = mysql.First(&packet, packetid)
	if result.RowsAffected == 0 {
		log.Errorf("Invalid envelop id: %d, block it.", packetid)
		c.JSON(200, gin.H{"code": OPEN_INVALID_PACKET, "msg": OPEN_INVALID_PACKET_MESSAGE, "data": gin.H{}})
		return
	}

	if packet.Opened {
		log.Errorf("Envelop %d has been opened yet.", packetid)
		c.JSON(200, gin.H{"code": OPEN_REPEAT, "msg": OPEN_REPEAT_MESSAGE, "data": gin.H{}})
		return
	}

	if userid != packet.UserID {
		log.Errorf("User %d don't own envelop %d", userid, packetid)
		c.JSON(200, gin.H{"code": OPEN_NOT_MATCH, "msg": OPEN_NOT_MATCH_MESSAGE, "data": gin.H{}})
		return
	}

	user.Balance += packet.Value
	packet.Opened = true
	mysql.Save(&user)
	mysql.Save(&packet)

	c.JSON(200, gin.H{"code": OPEN_SUCCESS, "msg": OPEN_SUCCESS_MESSAGE, "data": gin.H{"value": int32(packet.Value * 100)}})
}

func WalletListHandler(c *gin.Context) {
	uid := common.ConvertString(c.PostForm("uid"), "int32").(int32)
	log.Infof("Query %d's wallet", uid)

	user := db.User{UserID: uid, Amount: 0, Balance: 0.}
	mysql.Where(db.User{UserID: uid}).FirstOrCreate(&user)

	packets, _ := common.GetRedPacketsByUID(mysql, uid)
	envelops := []gin.H{}
	for _, p := range packets {
		envelops = append(envelops, p.JsonFormat())
	}

	c.JSON(200, gin.H{
		"code": WALLET_SUCCESS,
		"msg":  WALLET_SUCCESS_MESSAGE,
		"data": gin.H{
			"amount":       int32(user.Balance * 100),
			"envelop_list": envelops,
		},
	})
}
