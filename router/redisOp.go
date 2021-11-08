/*
 * @Author: your name
 * @Date: 2021-11-07 12:17:43
 * @LastEditTime: 2021-11-07 19:22:24
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
 * @description:
 * @param {*gin.Context} c
 * @param {int32} newTotalNum
 * @param {int64} newTotalmoney
 * @return {*} -1:error for lua 	0: envelop has been opened		1: change success
 */
func ChangeConfig(c *gin.Context, newTotalNum int32, newTotalmoney int64) int {
	changeScript := db.GenerateChangeScript()
	diffNum := server.sysconfig.TotalNum - newTotalNum
	diffMoney := server.sysconfig.TotalMoney - newTotalmoney
	ret, err := changeScript.Run(server.redisdb, []string{"RemainNum", "RemainMoney"}, diffNum, diffMoney).Result()
	if err != nil || ret.(int64) == -1 {
		log.Debug(err)
		c.JSON(http.StatusOK, gin.H{"code": constant.CHANGE_INVALID, "msg": constant.CHANGE_INVALID_MESSAGE, "data": gin.H{}})
		return -1
	}
	if ret.(int64) == 0 {
		log.Errorf("")
		c.JSON(http.StatusOK, gin.H{"code": constant.CHANGE_INVALID, "msg": constant.CHANGE_INVALID_MESSAGE, "data": gin.H{}})
		return 0
	}
	return 1
}
