package src

import (
	"math/rand"
	"time"
)

type RedPacket struct {
	value     float32
	opened    bool
	timestamp int64
	packetid  int64
}

func GetRedPacket(remain_num int, remain_money float32, min_money float32, max_money float32) RedPacket {
	if remain_num == 1 {
		return RedPacket{
			value:     min(remain_money, max_money),
			opened:    false,
			timestamp: time.Now().UnixNano(),
			packetid:  time.Now().UnixNano(),
		}
	}
	mean_money := remain_money / float32(remain_num)
	max_money = min(max_money, 2*mean_money-min_money)
	money := min_money + (max_money-min_money)*Rand()
	return RedPacket{
		value:     money,
		opened:    false,
		timestamp: time.Now().UnixNano(),
		packetid:  time.Now().UnixNano(),
	}
}

func Rand() float32 {
	return float32(rand.Intn(101)) / 100
}

func min(a float32, b float32) float32 {
	if a <= b {
		return a
	}
	return b
}
