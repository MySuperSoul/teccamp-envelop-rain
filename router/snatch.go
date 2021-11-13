/*
 * @Author: your name
 * @Date: 2021-11-08 16:49:31
 * @LastEditTime: 2021-11-08 17:04:50
 * @LastEditors: Please set LastEditors
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /teccamp-envelop-rain/router/snatch.go
 */
package router

import (
	"encoding/json"
	"envelop-rain/common"
	"envelop-rain/constant"
	db "envelop-rain/repository"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SnatchHandler(c *gin.Context) {
	json_str := make(map[string]int32)
	if err := c.BindJSON(&json_str); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.SNATCH_JSON_PARSE_ERROR, "msg": constant.SNATCH_JSON_PARSE_ERROR_MESSAGE, "data": gin.H{}})
		return
	}

	uid := json_str["uid"]
	uidStr := fmt.Sprintf("%d", uid)
	// first to judge whether has packet left
	if server.sendall {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_NO_RED_PACKET, "msg": constant.SNATCH_NO_RED_PACKET_MESSAGE, "data": gin.H{}})
		return
	}

	if server.bloomFilter.TestString(uidStr) {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_EXCEED_MAX_AMOUNT, "msg": constant.SNATCH_EXCEED_MAX_AMOUNT_MESSAGE, "data": gin.H{}})
		return
	}

	if common.Rand() > server.sysconfig.P {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_NOT_LUCKY, "msg": constant.SNATCH_NOT_LUCKY_MESSAGE, "data": gin.H{}})
		return
	}

	// generate packet id
	packetid := <-server.IDGenerator.IDs
	ret, _ := db.SnatchScript().Run(server.redisdb, []string{"CurrentNum"}, uidStr, packetid, server.sysconfig.TotalNum, server.sysconfig.MaxAmount, server.sysconfig.P).Result()
	retf := int(ret.(int64))
	if retf == constant.SNATCH_NO_RED_PACKET {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_NO_RED_PACKET, "msg": constant.SNATCH_NO_RED_PACKET_MESSAGE, "data": gin.H{}})
		server.sendall = true
		return
	}
	if retf == constant.SNATCH_EXCEED_MAX_AMOUNT {
		c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_EXCEED_MAX_AMOUNT, "msg": constant.SNATCH_EXCEED_MAX_AMOUNT_MESSAGE, "data": gin.H{}})
		server.bloomFilter.AddString(uidStr)
		return
	}

	if retf == int(server.sysconfig.MaxAmount) {
		server.bloomFilter.AddString(uidStr)
	}

	// success snatch
	c.JSON(http.StatusOK, gin.H{
		"code": constant.SNATCH_SUCCESS,
		"msg":  constant.SNATCH_SUCCESS_MESSAGE,
		"data": gin.H{"envelop_id": packetid, "max_count": server.sysconfig.MaxAmount, "cur_count": retf},
	})

	// insert the redpacket
	timestamp := time.Now().UnixNano()
	server.redisdb.HMSet("packet-"+fmt.Sprint(packetid), map[string]interface{}{
		"userid":    uid,
		"value":     0,
		"opened":    false,
		"timestamp": timestamp,
	})

	// send to database to create the redpacket
	packet_info := map[string]interface{}{
		"type":      constant.CREATE_PACKET_TYPE,
		"uid":       uid,
		"packet_id": packetid,
		"timestamp": timestamp,
	}
	message, _ := json.Marshal(packet_info)
	server.producer.SendDBMessage(message)
}
