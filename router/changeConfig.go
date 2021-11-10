/*
 * @Author: your name
 * @Date: 2021-11-08 16:49:31
 * @LastEditTime: 2021-11-08 16:56:21
 * @LastEditors: your name
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /teccamp-envelop-rain/router/change_config.go
 */
package router

import (
	"envelop-rain/constant"
	db "envelop-rain/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 * @description:
 * @param {*gin.Context} c
 * @param {int64} newTotalmoney
 * @return {*} -1:error for lua 	0: envelop has been opened		1: change success
 */
func ChangeConfig(c *gin.Context, newTotalmoney int64) int {
	changeScript := db.GenerateChangeScript()
	diffMoney := server.sysconfig.TotalMoney - newTotalmoney
	ret, err := changeScript.Run(server.redisdb, []string{"RemainMoney"}, diffMoney).Result()
	if err != nil || ret.(int64) == -1 {
		server.logger.Debug(err)
		c.JSON(http.StatusOK, gin.H{"code": constant.CHANGE_INVALID, "msg": constant.CHANGE_INVALID_MESSAGE, "data": gin.H{}})
		return -1
	}
	return 1
}

type money_setting struct {
	NewTotalmoney int64 `json:"totalmoney"`
}

func ChangeConfigsHandler(c *gin.Context) {
	var setting money_setting
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": constant.CHANGE_JSON_PARSE_ERROR, "msg": constant.CHANGE_JSON_PARSE_ERROR_MESSAGE})
		return
	}
	//update the config
	ret := ChangeConfig(c, setting.NewTotalmoney)
	if ret < 1 {
		return
	}
	server.sysconfig.TotalMoney = setting.NewTotalmoney
	message := fmt.Sprintf("TotalMoney: from %d to %d", server.sysconfig.TotalMoney, setting.NewTotalmoney)
	server.logger.Debug(message)
	c.JSON(http.StatusOK, gin.H{"code": constant.CHANGE_SUCCESS, "msg": constant.CHANGE_SUCCESS_MESSAGE})
}
