package router

import (
	"envelop-rain/configs"
	db "envelop-rain/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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
	uid := fmt.Sprint(json_str["uid"])
	log.Infof("Query %s's wallet", uid)

	packets, balance := db.GetRedPacketsByUID(server.redisdb, uid)
	envelops := []gin.H{}
	for _, p := range packets {
		envelops = append(envelops, p.JsonFormat())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": WALLET_SUCCESS,
		"msg":  WALLET_SUCCESS_MESSAGE,
		"data": gin.H{
			"amount":       balance,
			"envelop_list": envelops,
		},
	})
}

// Only for testing in cloud service
func FlushDBHandler(c *gin.Context) {
	server.redisdb.FlushDB()
	configs.SetConfigToRedis(&server.sysconfig, server.redisdb)
	db.GenerateTables(server.mysql)
}
