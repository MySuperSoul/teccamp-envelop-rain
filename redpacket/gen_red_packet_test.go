package redpacket

import (
	"envelop-rain/common"
	"testing"
)

func TestGetRedPacket(t *testing.T) {
	common.SetRandomSeed()
	cases := []struct {
		total_num, remain_num                           int
		total_money, remain_money, min_money, max_money float32
	}{
		{100, 100, 10000., 10000., 1., 200.},
		{128, 128, 23456, 23456, 10., 300.},
		{123456, 123456, 12345678., 12345678., 5., 400.},
	}

	for _, c := range cases {
		// iterate each red packet to check
		for i := 0; i < c.total_num; i++ {
			red_packet := GetRedPacket(int32(c.remain_num), c.remain_money, c.min_money, c.max_money)
			if red_packet.Value < c.min_money || red_packet.Value > c.max_money {
				t.Fatalf("Get error on packet value, %f is not in range %f -> %f", red_packet.Value, c.min_money, c.max_money)
			}
			c.remain_num--
			c.remain_money -= red_packet.Value
		}

		// check the final remain money, can not exceed the total money
		if c.remain_money < 0 {
			t.Fatalf("Final Remain money error, value is %f", c.remain_money)
		}
	}
}
