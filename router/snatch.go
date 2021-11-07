package router

import (
	"envelop-rain/common"
	db "envelop-rain/repository"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	uidStr := fmt.Sprintf("%d", uid)
	// first to judge whether has packet left
	current_num := db.GetSingleValueFromRedis(server.redisdb, "CurrentNum", "int32").(int32)

	if current_num == server.sysconfig.TotalNum {
		c.JSON(http.StatusOK, gin.H{"code": SNATCH_NO_RED_PACKET, "msg": SNATCH_NO_RED_PACKET_MESSAGE, "data": gin.H{}})
		return
	}
	// Then perform later operations
	// First judge whether has this user
	if n, _ := server.redisdb.Exists(uidStr).Result(); n == 0 { // no this user
		server.redisdb.SetNX("user-"+uidStr, 0, 0) // set amount to 0
		// TODO: send to database to create this user with balance = 0
	}

	// Check whether exceed max amount
	amount, _ := server.redisdb.Get("user-" + uidStr).Int64()
	if int32(amount) == server.sysconfig.MaxAmount {
		c.JSON(http.StatusOK, gin.H{"code": SNATCH_EXCEED_MAX_AMOUNT, "msg": SNATCH_EXCEED_MAX_AMOUNT_MESSAGE, "data": gin.H{}})
		return
	}

	// Then to judge whether the user is lucky enough
	if common.Rand() > server.sysconfig.P {
		c.JSON(http.StatusOK, gin.H{"code": SNATCH_NOT_LUCKY, "msg": SNATCH_NOT_LUCKY_MESSAGE, "data": gin.H{}})
		return
	}

	// generate packet_id
	packetid := time.Now().UnixNano() / 1000
	// send message
	c.JSON(http.StatusOK, gin.H{
		"code": SNATCH_SUCCESS,
		"msg":  SNATCH_SUCCESS_MESSAGE,
		"data": gin.H{"envelop_id": packetid, "max_count": server.sysconfig.MaxAmount, "cur_count": int32(amount) + 1},
	})

	// update user amount
	server.redisdb.Incr("CurrentNum")
	server.redisdb.Incr("user-" + uidStr)
	// insert the redpacket
	server.redisdb.HMSet("packet-"+fmt.Sprint(packetid), map[string]interface{}{
		"userid":    uid,
		"value":     0,
		"opened":    false,
		"timestamp": time.Now().UnixNano(),
	})
	server.redisdb.LPush(uidStr+"-wallet", packetid)

	// TODO: send to database to create the redpacket
}
