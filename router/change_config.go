package router

import (
	"envelop-rain/constant"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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
