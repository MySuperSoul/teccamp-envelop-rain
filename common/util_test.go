package common

import (
	"testing"
)

func TestRand(t *testing.T) {
	loop_times := 100
	for i := 0; i < loop_times; i++ {
		ratio := Rand()
		if ratio < 0 || ratio > 1 {
			t.Fatalf("Out of boundary, value is %f", ratio)
		}
	}
}

func TestGetRedPacket(t *testing.T) {
	SetRandomSeed()
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
			if red_packet < c.min_money || red_packet > c.max_money {
				t.Fatalf("Get error on packet value, %d is not in range %d -> %d", red_packet, c.min_money, c.max_money)
			}
			c.remain_num--
			c.remain_money -= int64(red_packet)
		}

		// check the final remain money, can not exceed the total money
		if c.remain_money < 0 {
			t.Fatalf("Final Remain money error, value is %d", c.remain_money)
		}
	}
}
