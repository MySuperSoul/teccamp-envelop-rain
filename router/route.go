/*
 * @Author: your name
 * @Date: 2021-11-02 19:16:51
 * @LastEditTime: 2021-11-03 21:26:08
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/router/route.go
 */
package router

import (
	"envelop-rain/common"
	"envelop-rain/controller"
	db "envelop-rain/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SnatchHandler(c *gin.Context) {
	json_str := make(map[string]int32)
	if err := c.BindJSON(&json_str); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": SNATCH_JSON_PARSE_ERROR, "msg": SNATCH_JSON_PARSE_ERROR_MESSAGE, "data": gin.H{}})
		return
	}
	if _, ok := json_str["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": SNATCH_EMPTY_UID, "msg": SNATCH_EMPTY_UID_MESSAGE, "data": gin.H{}})
		return
	}
	uid := json_str["uid"]
	log.Infof("snatch by user: %d", uid)

	// first to judge whether has packet left
	remain_num := db.GetSingleValueFromRedis(server.redisdb, "RemainNum", "int32").(int32)
	remain_money := db.GetSingleValueFromRedis(server.redisdb, "RemainMoney", "int64").(int64)

	if remain_num == 0 {
		c.JSON(http.StatusOK, gin.H{"code": SNATCH_NO_RED_PACKET, "msg": SNATCH_NO_RED_PACKET_MESSAGE, "data": gin.H{}})
		return
	}

	// Then perform later operations
	// get user information
	user := db.User{UserID: uid, Amount: 0, Balance: 0.}
	server.mysql.Where(db.User{UserID: uid}).FirstOrCreate(&user)

	// Then to check the maxamount
	max_amount := db.GetSingleValueFromRedis(server.redisdb, "MaxAmount", "int32").(int32)
	if user.Amount == max_amount {
		c.JSON(http.StatusOK, gin.H{"code": SNATCH_EXCEED_MAX_AMOUNT, "msg": SNATCH_EXCEED_MAX_AMOUNT_MESSAGE, "data": gin.H{}})
		return
	}

	// Then to judge whether the user is lucky enough
	if p := db.GetSingleValueFromRedis(server.redisdb, "P", "float32").(float32); common.Rand() > p {
		c.JSON(http.StatusOK, gin.H{"code": SNATCH_NOT_LUCKY, "msg": SNATCH_NOT_LUCKY_MESSAGE, "data": gin.H{}})
		return
	}

	// First generate the red packet
	packet := controller.GetRedPacket(remain_num, remain_money, server.sysconfig.MinMoney, server.sysconfig.MaxMoney)
	packet.UserID = uid

	// update remain value to redis
	server.redisdb.Decr("RemainNum")
	server.redisdb.Set("RemainMoney", remain_money-int64(packet.Value), 0)

	// update user amount and insert the red packet
	user.Amount++
	server.mysql.Model(&user).Update("amount", user.Amount)
	server.mysql.Create(&packet)

	// send message
	c.JSON(http.StatusOK, gin.H{
		"code": SNATCH_SUCCESS,
		"msg":  SNATCH_SUCCESS_MESSAGE,
		"data": gin.H{"envelop_id": packet.PacketID, "max_count": max_amount, "cur_count": user.Amount},
	})
}

type uid_envelopid struct {
	Uid      int32 `json:"uid"`
	Packetid int64 `json:"envelop_id"`
}

func OpenHandler(c *gin.Context) {
	var pair_id uid_envelopid
	if err := c.ShouldBindJSON(&pair_id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": OPEN_JSON_PARSE_ERROR, "msg": OPEN_JSON_PARSE_ERROR_MESSAGE, "data": gin.H{}})
		return
	}

	userid := pair_id.Uid
	packetid := pair_id.Packetid
	log.Infof("Envelop %d opened by %d.", packetid, userid)

	var user db.User
	var packet db.RedPacket
	result := server.mysql.First(&user, userid)
	if result.RowsAffected == 0 {
		log.Errorf("Invalid user id: %d, block him.", userid)
		c.JSON(http.StatusOK, gin.H{"code": OPEN_INVALID_USER, "msg": OPEN_INVALID_USER_MESSAGE, "data": gin.H{}})
		return
	}
	result = server.mysql.First(&packet, packetid)
	if result.RowsAffected == 0 {
		log.Errorf("Invalid envelop id: %d, block it.", packetid)
		c.JSON(http.StatusOK, gin.H{"code": OPEN_INVALID_PACKET, "msg": OPEN_INVALID_PACKET_MESSAGE, "data": gin.H{}})
		return
	}

	if packet.Opened {
		log.Errorf("Envelop %d has been opened yet.", packetid)
		c.JSON(http.StatusOK, gin.H{"code": OPEN_REPEAT, "msg": OPEN_REPEAT_MESSAGE, "data": gin.H{}})
		return
	}

	if userid != packet.UserID {
		log.Errorf("User %d don't own envelop %d", userid, packetid)
		c.JSON(http.StatusOK, gin.H{"code": OPEN_NOT_MATCH, "msg": OPEN_NOT_MATCH_MESSAGE, "data": gin.H{}})
		return
	}

	user.Balance += packet.Value
	packet.Opened = true
	server.mysql.Save(&user)
	server.mysql.Save(&packet)

	c.JSON(http.StatusOK, gin.H{"code": OPEN_SUCCESS, "msg": OPEN_SUCCESS_MESSAGE, "data": gin.H{"value": packet.Value}})
}

func WalletListHandler(c *gin.Context) {
	json_str := make(map[string]int32)
	if err := c.BindJSON(&json_str); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": WALLET_JSON_PARSE_ERROR, "msg": WALLET_JSON_PARSE_ERROR_MESSAGE, "data": gin.H{}})
		return
	}
	if _, ok := json_str["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": WALLET_EMPTY_ID, "msg": WALLET_EMPTY_ID_MESSAGE, "data": gin.H{}})
		return
	}
	uid := json_str["uid"]
	log.Infof("Query %d's wallet", uid)

	user := db.User{UserID: uid, Amount: 0, Balance: 0.}
	server.mysql.Where(db.User{UserID: uid}).FirstOrCreate(&user)

	packets, _ := db.GetRedPacketsByUID(server.mysql, uid)
	envelops := []gin.H{}
	for _, p := range packets {
		envelops = append(envelops, p.JsonFormat())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": WALLET_SUCCESS,
		"msg":  WALLET_SUCCESS_MESSAGE,
		"data": gin.H{
			"amount":       user.Balance,
			"envelop_list": envelops,
		},
	})
}
