package router

import (
	"envelop-rain/constant"
	db "envelop-rain/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WalletListHandler(c *gin.Context) {
	json_str := make(map[string]int32)
	if err := c.BindJSON(&json_str); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.WALLET_JSON_PARSE_ERROR, "msg": constant.WALLET_JSON_PARSE_ERROR_MESSAGE, "data": gin.H{}})
		return
	}

	uid := fmt.Sprint(json_str["uid"])

	packets, balance := db.GetRedPacketsByUID(server.redisdb, server.mysql, uid)
	envelops := []gin.H{}
	for _, p := range packets {
		envelops = append(envelops, p.JsonFormat())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": constant.WALLET_SUCCESS,
		"msg":  constant.WALLET_SUCCESS_MESSAGE,
		"data": gin.H{
			"amount":       balance,
			"envelop_list": envelops,
		},
	})
}
