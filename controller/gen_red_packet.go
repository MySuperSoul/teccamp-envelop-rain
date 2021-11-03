/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-02 21:04:36
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/controller/gen_red_packet.go
 */
package controller

import (
	"envelop-rain/common"
	db "envelop-rain/repository"
	"math/rand"
	"time"
)

func GetRedPacket(remain_num int32, remain_money int64, min_money int32, max_money int32) db.RedPacket {
	if remain_num == 1 {
		return db.RedPacket{
			Value:     common.GetMin(int32(remain_money), max_money),
			Opened:    false,
			Timestamp: time.Now().UnixNano(),
			PacketID:  time.Now().UnixNano(),
		}
	}
	mean_money := int32(remain_money / int64(remain_num))
	max_money = common.GetMin(max_money, 2*mean_money-min_money)
	money := min_money + rand.Int31n(max_money-min_money+1)
	return db.RedPacket{
		Value:     money,
		Opened:    false,
		Timestamp: time.Now().UnixNano(),
		PacketID:  time.Now().UnixNano(),
	}
}
