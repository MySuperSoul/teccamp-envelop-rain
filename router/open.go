package router

import (
	"encoding/json"
	"envelop-rain/common"
	"fmt"
	"net/http"
	"time"

	"envelop-rain/constant"
	db "envelop-rain/repository"

	"github.com/gin-gonic/gin"
)

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

	if n, _ := server.redisdb.Exists("packet-" + packetid).Result(); n == 0 {
		// Then check db
		var packet db.RedPacket
		if result := server.mysql.Where("packet_id = ?", pair_id.Packetid).Find(&packet); result.RowsAffected == 0 {
			server.logger.Errorf("Invalid envelop id: %s, block it.", packetid)
			c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_INVALID_PACKET, "msg": constant.OPEN_INVALID_PACKET_MESSAGE, "data": gin.H{}})
			return
		} else { // set to redis
			server.redisdb.HMSet("packet-"+packetid, map[string]interface{}{
				"userid":    packet.UserID,
				"value":     packet.Value,
				"opened":    packet.Opened,
				"timestamp": packet.Timestamp,
			})
		}
	}

	if puid, _ := server.redisdb.HGet("packet-"+packetid, "userid").Result(); puid != uid {
		server.logger.Errorf("User %s don't own envelop %s", uid, packetid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_NOT_MATCH, "msg": constant.OPEN_NOT_MATCH_MESSAGE, "data": gin.H{}})
		return
	}

	if isopen, _ := server.redisdb.HGet("packet-"+packetid, "opened").Result(); common.ConvertString(isopen, "bool").(bool) {
		server.logger.Errorf("Envelop %s has been opened yet.", packetid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_REPEAT, "msg": constant.OPEN_REPEAT_MESSAGE, "data": gin.H{}})
		return
	}

	if server.nomoney {
		server.logger.Error("Money has send out")
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_NO_MONEY, "msg": constant.OPEN_NO_MONEY_MESSAGE, "data": gin.H{}})
		return
	}

	// generate value here
	v, _ := db.GeneratePacketScript().Run(server.redisdb, []string{"RemainNum", "RemainMoney"}, server.sysconfig.MaxMoney, server.sysconfig.MinMoney, time.Now().Nanosecond()).Result()
	value := int32(v.(int64))

	if value == -1 {
		server.nomoney = true
		server.logger.Error("Money has send out")
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_NO_MONEY, "msg": constant.OPEN_NO_MONEY_MESSAGE, "data": gin.H{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_SUCCESS, "msg": constant.OPEN_SUCCESS_MESSAGE, "data": gin.H{"value": value}})

	server.redisdb.HMSet("packet-"+packetid, map[string]interface{}{
		"value":  value,
		"opened": true,
	})

	// Update opened field and value to packet table and update remain
	packet_info := map[string]interface{}{
		"type":      constant.UPDATE_PACKET_TYPE,
		"packet_id": pair_id.Packetid,
		"value":     value,
	}
	data, _ := json.Marshal(packet_info)
	server.producer.SendDBMessage(data)

	info := map[string]interface{}{
		"type":  constant.UPDATE_REMAIN_TYPE,
		"money": value,
	}
	info_data, _ := json.Marshal(info)
	server.producer.SendDBMessage(info_data)
}
