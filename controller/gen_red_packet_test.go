/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-02 21:02:39
 * @LastEditors: your name
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/controller/gen_red_packet_test.go
 */
package controller

import (
	"envelop-rain/common"
	"testing"
)

func TestGetRedPacket(t *testing.T) {
	common.SetRandomSeed()
	cases := []struct {
		total_num, remain_num     int
		total_money, remain_money int64
		min_money, max_money      int32
	}{
		{100, 100, 10000, 10000, 1, 200},
		{128, 128, 23456, 23456, 10, 300},
		{123456, 123456, 12345678, 12345678, 5, 400},
	}

	for _, c := range cases {
		// iterate each red packet to check
		for i := 0; i < c.total_num; i++ {
			red_packet := GetRedPacket(int32(c.remain_num), c.remain_money, c.min_money, c.max_money)
			if red_packet.Value < c.min_money || red_packet.Value > c.max_money {
				t.Fatalf("Get error on packet value, %d is not in range %d -> %d", red_packet.Value, c.min_money, c.max_money)
			}
			c.remain_num--
			c.remain_money -= int64(red_packet.Value)
		}

		// check the final remain money, can not exceed the total money
		if c.remain_money < 0 {
			t.Fatalf("Final Remain money error, value is %d", c.remain_money)
		}
	}
}
