package redpacket

import (
	"envelop-rain/common"
	"time"
)

func GetRedPacket(remain_num int, remain_money float32, min_money float32, max_money float32) RedPacket {
	if remain_num == 1 {
		return RedPacket{
			value:     common.GetMin(remain_money, max_money),
			opened:    false,
			timestamp: time.Now().UnixNano(),
			packetid:  time.Now().UnixNano(),
		}
	}
	mean_money := remain_money / float32(remain_num)
	max_money = common.GetMin(max_money, 2*mean_money-min_money)
	money := min_money + (max_money-min_money)*common.Rand()
	return RedPacket{
		value:     money,
		opened:    false,
		timestamp: time.Now().UnixNano(),
		packetid:  time.Now().UnixNano(),
	}
}
