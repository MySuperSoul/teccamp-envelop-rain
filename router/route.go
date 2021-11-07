/*
 * @Author: your name
 * @Date: 2021-11-02 19:16:51
 * @LastEditTime: 2021-11-07 17:34:05
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/router/route.go
 */
package router

import (
	"envelop-rain/common"
	"envelop-rain/constant"
	"envelop-rain/controller"
	db "envelop-rain/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SnatchHandler(c *gin.Context) {
	json_str := make(map[string]int32)
	if err := c.BindJSON(&json_str); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_JSON_PARSE_ERROR, "msg": constant.SNATCH_JSON_PARSE_ERROR_MESSAGE, "data": gin.H{}})
		return
	}
	if _, ok := json_str["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_EMPTY_UID, "msg": constant.SNATCH_EMPTY_UID_MESSAGE, "data": gin.H{}})
		return
	}
	uid := json_str["uid"]
	uidStr := fmt.Sprintf("%d", uid)
	log.Infof("snatch by user: %d", uid)
	// first to judge whether has packet left
	remain_num := db.GetSingleValueFromRedis(server.redisdb, "RemainNum", "int32").(int32)

	if remain_num == 0 {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_NO_RED_PACKET, "msg": constant.SNATCH_NO_RED_PACKET_MESSAGE, "data": gin.H{}})
		return
	}
	// Then perform later operations
	// First judge whether has this user
	if n, _ := server.redisdb.Exists(uidStr).Result(); n == 0 { // no this user
		server.redisdb.HMSet(uidStr, map[string]interface{}{"amount": 0, "balance": 0})
		// TODO: send to database to create this user with balance = 0
	}

	// Check whether exceed max amount
	if amount, _ := server.redisdb.HGet(uidStr, "amount").Int64(); int32(amount) == server.sysconfig.MaxAmount {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_EXCEED_MAX_AMOUNT, "msg": constant.SNATCH_EXCEED_MAX_AMOUNT_MESSAGE, "data": gin.H{}})
		return
	}

	// Then to judge whether the user is lucky enough
	if common.Rand() > server.sysconfig.P {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_NOT_LUCKY, "msg": constant.SNATCH_NOT_LUCKY_MESSAGE, "data": gin.H{}})
		return
	}

	remain_money := db.GetSingleValueFromRedis(server.redisdb, "RemainMoney", "int64").(int64)
	// First generate the red packet
	packet := controller.GetRedPacket(remain_num, remain_money, server.sysconfig.MinMoney, server.sysconfig.MaxMoney)
	packet.UserID = uid

	// update remain value to redis
	if ret := UpdateRemain(c, packet.Value); ret < 1 {
		return
	}

	// update user amount
	cur_count := UpdateUserAmount(c, uidStr)
	if cur_count < 1 {
		return
	}

	// insert the redpacket
	server.redisdb.HMSet(fmt.Sprint(packet.PacketID), packet.ToRedisFormat())
	server.redisdb.LPush(uidStr+"-wallet", packet.PacketID)

	// TODO: send to database to create the redpacket

	// send message
	c.JSON(http.StatusOK, gin.H{
		"code": constant.SNATCH_SUCCESS,
		"msg":  constant.SNATCH_SUCCESS_MESSAGE,
		"data": gin.H{"envelop_id": packet.PacketID, "max_count": server.sysconfig.MaxAmount, "cur_count": cur_count},
	})
}

type uid_envelopid struct {
	Uid      int32 `json:"uid"`
	Packetid int64 `json:"envelop_id"`
}

func OpenHandler(c *gin.Context) {
	var pair_id uid_envelopid
	if err := c.ShouldBindJSON(&pair_id); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_JSON_PARSE_ERROR, "msg": constant.OPEN_JSON_PARSE_ERROR_MESSAGE, "data": gin.H{}})
		return
	}

	uid := fmt.Sprint(pair_id.Uid)
	packetid := fmt.Sprint(pair_id.Packetid)
	log.Infof("Envelop %s opened by %s.", packetid, uid)

	// invalid user here
	if n, _ := server.redisdb.Exists(uid).Result(); n == 0 {
		log.Errorf("Invalid user id: %s, block him.", uid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_INVALID_USER, "msg": constant.OPEN_INVALID_USER_MESSAGE, "data": gin.H{}})
		return
	}

	if n, _ := server.redisdb.Exists(packetid).Result(); n == 0 {
		log.Errorf("Invalid envelop id: %s, block it.", packetid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_INVALID_PACKET, "msg": constant.OPEN_INVALID_PACKET_MESSAGE, "data": gin.H{}})
		return
	}

	if puid, _ := server.redisdb.HGet(packetid, "userid").Result(); puid != uid {
		log.Errorf("User %s don't own envelop %s", uid, packetid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_NOT_MATCH, "msg": constant.OPEN_NOT_MATCH_MESSAGE, "data": gin.H{}})
		return
	}

	//try open the red packet
	ret := OpenPacket(c, uid, packetid)
	if ret < 1 {
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_SUCCESS, "msg": constant.OPEN_SUCCESS_MESSAGE, "data": gin.H{"value": int32(ret)}})
}

func WalletListHandler(c *gin.Context) {
	json_str := make(map[string]int32)
	if err := c.BindJSON(&json_str); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": constant.WALLET_JSON_PARSE_ERROR, "msg": constant.WALLET_JSON_PARSE_ERROR_MESSAGE, "data": gin.H{}})
		return
	}
	if _, ok := json_str["uid"]; !ok {
		c.JSON(http.StatusOK, gin.H{"code": constant.WALLET_EMPTY_ID, "msg": constant.WALLET_EMPTY_ID_MESSAGE, "data": gin.H{}})
		return
	}
	uid := fmt.Sprint(json_str["uid"])
	log.Infof("Query %s's wallet", uid)

	packets, _ := db.GetRedPacketsByUID(server.redisdb, uid)
	envelops := []gin.H{}
	for _, p := range packets {
		envelops = append(envelops, p.JsonFormat())
	}

	balance, _ := server.redisdb.HGet(uid, "balance").Int64()

	c.JSON(http.StatusOK, gin.H{
		"code": constant.WALLET_SUCCESS,
		"msg":  constant.WALLET_SUCCESS_MESSAGE,
		"data": gin.H{
			"amount":       int32(balance),
			"envelop_list": envelops,
		},
	})
}

type money_setting struct {
	NewTotalmoney int64 `json:"totalmoney"`
	NewTotalNum   int32 `json:"totalnum"`
}

func ChangeConfigsHandler(c *gin.Context) {
	var setting money_setting
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": constant.CHANGE_JSON_PARSE_ERROR, "msg": constant.CHANGE_JSON_PARSE_ERROR_MESSAGE})
		return
	}
	//update the config
	ret := ChangeConfig(c, setting.NewTotalNum, setting.NewTotalmoney)
	if ret < 1 {
		return
	}
	server.sysconfig.TotalMoney = setting.NewTotalmoney
	server.sysconfig.TotalNum = setting.NewTotalNum
	message := fmt.Sprintf("TotalMoney: from %d to %d	TotalNum: from %d to %d",
		server.sysconfig.TotalMoney, setting.NewTotalmoney,
		server.sysconfig.TotalNum, setting.NewTotalNum)
	log.Debug(message)
	c.JSON(http.StatusOK, gin.H{"code": constant.CHANGE_SUCCESS, "msg": constant.CHANGE_SUCCESS_MESSAGE})

}
