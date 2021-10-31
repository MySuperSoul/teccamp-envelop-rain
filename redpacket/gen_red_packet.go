package redpacket

import (
	"envelop-rain/common"
	"envelop-rain/db"
	"time"
)

func GetRedPacket(remain_num int32, remain_money float32, min_money float32, max_money float32) db.RedPacket {
	if remain_num == 1 {
		return db.RedPacket{
			Value:     common.GetMin(remain_money, max_money),
			Opened:    false,
			Timestamp: time.Now().UnixNano(),
			PacketID:  time.Now().UnixNano(),
		}
	}
	mean_money := remain_money / float32(remain_num)
	max_money = common.GetMin(max_money, 2*mean_money-min_money)
	money := min_money + (max_money-min_money)*common.Rand()
	return db.RedPacket{
		Value:     money,
		Opened:    false,
		Timestamp: time.Now().UnixNano(),
		PacketID:  time.Now().UnixNano(),
	}
}
