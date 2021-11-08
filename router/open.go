package router

import (
	"envelop-rain/common"
	"fmt"
	"net/http"
	"time"

	"envelop-rain/constant"
	db "envelop-rain/repository"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	log.Infof("Envelop %s opened by %s.", packetid, uid)

	// invalid user here
	if n, _ := server.redisdb.Exists(uid + "-wallet").Result(); n == 0 {
		log.Errorf("Invalid user id: %s, block him.", uid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_INVALID_USER, "msg": constant.OPEN_INVALID_USER_MESSAGE, "data": gin.H{}})
		return
	}

	if n, _ := server.redisdb.Exists("packet-" + packetid).Result(); n == 0 {
		log.Errorf("Invalid envelop id: %s, block it.", packetid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_INVALID_PACKET, "msg": constant.OPEN_INVALID_PACKET_MESSAGE, "data": gin.H{}})
		return
	}

	if puid, _ := server.redisdb.HGet("packet-"+packetid, "userid").Result(); puid != uid {
		log.Errorf("User %s don't own envelop %s", uid, packetid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_NOT_MATCH, "msg": constant.OPEN_NOT_MATCH_MESSAGE, "data": gin.H{}})
		return
	}

	if isopen, _ := server.redisdb.HGet("packet-"+packetid, "opened").Result(); common.ConvertString(isopen, "bool").(bool) {
		log.Errorf("Envelop %s has been opened yet.", packetid)
		c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_REPEAT, "msg": constant.OPEN_REPEAT_MESSAGE, "data": gin.H{}})
		return
	}

	// generate value here
	v, _ := db.GeneratePacketScript().Run(server.redisdb, []string{"RemainNum", "RemainMoney"}, server.sysconfig.MaxMoney, server.sysconfig.MinMoney, time.Now().Nanosecond()).Result()
	value := int32(v.(int64))
	c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_SUCCESS, "msg": constant.OPEN_SUCCESS_MESSAGE, "data": gin.H{"value": value}})

	server.redisdb.HMSet("packet-"+packetid, map[string]interface{}{
		"value":  value,
		"opened": true,
	})

	// TODO: Update balance to user table
	// TODO: Update opened field to packet table
}
