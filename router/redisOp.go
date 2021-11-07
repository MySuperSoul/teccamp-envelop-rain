/*
 * @Author: your name
 * @Date: 2021-11-07 12:17:43
 * @LastEditTime: 2021-11-07 15:29:25
 * @LastEditors: Please set LastEditors
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /teccamp-envelop-rain/router/redisOp.go
 */
package router

import (
	"envelop-rain/constant"
	db "envelop-rain/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

/**
 * @description: update the remain money and remain number for red packet
 * @param {*gin.Context} c
 * @param {int32} the value for the moeny in packet
 * @return {*} -1:error for lua 	0:no red packet remain 		1:snatch red packet successfully
 */
func UpdateRemain(c *gin.Context, value int32) int {
	decrScript := db.GenerateDecrScript()
	ret, err := decrScript.Run(server.redisdb, []string{"RemainNum", "RemainMoney"}, value).Result()
	if err != nil || ret.(int64) == -1 {
		log.Debug(err)
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_NOT_LUCKY, "msg": constant.SNATCH_NOT_LUCKY_MESSAGE, "data": gin.H{}})
		return -1
	}
	if ret.(int64) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_NO_RED_PACKET, "msg": constant.SNATCH_NO_RED_PACKET_MESSAGE, "data": gin.H{}})
		return 0
	}
	return 1
}

/**
 * @description: update the user amount
 * @param {*gin.Context} c
 * @param {string} uidStr
 * @return {*} -1:error for lua 	0:exceed the max amount 	others:success and return the user amount
 */
func UpdateUserAmount(c *gin.Context, uidStr string) int {
	amountDecrScript := db.GenerateUserAmountDecrScript()
	ret, err := amountDecrScript.Run(server.redisdb, []string{uidStr, "amount"}, server.sysconfig.MaxAmount).Result()
	if err != nil || ret.(int64) == -1 {
		log.Debug(err)
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_NOT_LUCKY, "msg": constant.SNATCH_NOT_LUCKY_MESSAGE, "data": gin.H{}})
		return -1
	}
	if ret.(int64) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_EXCEED_MAX_AMOUNT, "msg": constant.SNATCH_EXCEED_MAX_AMOUNT_MESSAGE, "data": gin.H{}})
		return 0
	}
	return int(ret.(int64))
}

/**
 * @description: try open the red packet
 * @param {*gin.Context} c
 * @param {string} uid
 * @param {string} packetid
 * @return {*} -1:error for lua 	0: envelop has been opened		others:success and return the user balance
 */
func OpenPacket(c *gin.Context, uid string, packetid string) int {
	openPacketScript := db.GenerateOpenPacketScript()
	ret, err := openPacketScript.Run(server.redisdb, []string{packetid, "opened", "value", uid, "balance"}).Result()
	if err != nil || ret.(int64) == -1 {
		log.Debug(err)
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_NOT_LUCKY, "msg": constant.SNATCH_NOT_LUCKY_MESSAGE, "data": gin.H{}})
		return -1
	}
	if ret.(int64) == 0 {
		log.Errorf("Envelop %s has been opened yet.", packetid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_REPEAT, "msg": constant.OPEN_REPEAT_MESSAGE, "data": gin.H{}})
		return 0
	}
	return int(ret.(int64))
}
