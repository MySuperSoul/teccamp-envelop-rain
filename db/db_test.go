/*
 * @Author: your name
 * @Date: 2021-11-01 14:56:45
 * @LastEditTime: 2021-11-01 15:57:24
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/db/db_test.go
 */
package db

import (
	"envelop-rain/configs"
	"testing"
)

func TestGetSingleValueFromRedis(t *testing.T) {
	redisdb := GetRedisClient()
	defer redisdb.Close()

	// Test int get
	val := 1234
	redisdb.Set("test", val, 0)
	if gval := GetSingleValueFromRedis(redisdb, "test", "int").(int); gval != val {
		t.Fatalf("Get int fail, target: %d, get: %d", val, gval)
	}
	// Test int32 get
	ival := int32(1234567890)
	redisdb.Set("test", ival, 0)
	if gval := GetSingleValueFromRedis(redisdb, "test", "int32").(int32); gval != ival {
		t.Fatalf("Get int fail, target: %d, get: %d", ival, gval)
	}
	// Test float32 get
	val2 := float32(1234.5678)
	redisdb.Set("test", val2, 0)
	if gval := GetSingleValueFromRedis(redisdb, "test", "float32").(float32); gval != val2 {
		t.Fatalf("Get float32 fail, target: %f, get: %f", val2, gval)
	}
	// Test float64 get
	val3 := 12345.6789
	redisdb.Set("test", val3, 0)
	if gval := GetSingleValueFromRedis(redisdb, "test", "float64").(float64); gval != val3 {
		t.Fatalf("Get float64 fail, target: %f, get: %f", val3, gval)
	}
	// Test string get
	val4 := "teccamp"
	redisdb.Set("test", val4, 0)
	if gval := GetSingleValueFromRedis(redisdb, "test", "string").(string); gval != val4 {
		t.Fatalf("Get string fail, target: %s, get: %s", val4, gval)
	}
}

func TestSetConfigToRedis(t *testing.T) {
	redisdb := GetRedisClient()
	defer redisdb.Close()

	cases := []configs.SystemConfig{
		{TotalMoney: 1., MinMoney: 1., MaxMoney: 1., P: 1., TotalNum: 1, MaxAmount: 1},
		{TotalMoney: 1., MinMoney: 10., MaxMoney: 100., P: 1000., TotalNum: 1234, MaxAmount: 4},
		{TotalMoney: 3.4, MinMoney: 10.5, MaxMoney: 34.5, P: 12.3, TotalNum: 9, MaxAmount: 8}}

	for _, c := range cases {
		// first set to redis
		configs.SetConfigToRedis(&c, redisdb)

		if val := GetSingleValueFromRedis(redisdb, "TotalMoney", "float32").(float32); val != c.TotalMoney {
			t.Fatalf("Restore from redis *total money* fail. Target: %f, Get: %f", c.TotalMoney, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "TotalNum", "int32").(int32); val != c.TotalNum {
			t.Fatalf("Restore from redis *total num* fail. Target: %d, Get: %d", c.TotalNum, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "RemainMoney", "float32").(float32); val != c.TotalMoney {
			t.Fatalf("Restore from redis *remain money* fail. Target: %f, Get: %f", c.TotalMoney, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "RemainNum", "int32").(int32); val != c.TotalNum {
			t.Fatalf("Restore from redis *remain num* fail. Target: %d, Get: %d", c.TotalNum, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "MinMoney", "float32").(float32); val != c.MinMoney {
			t.Fatalf("Restore from redis *min money* fail. Target: %f, Get: %f", c.MinMoney, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "MaxMoney", "float32").(float32); val != c.MaxMoney {
			t.Fatalf("Restore from redis *max money* fail. Target: %f, Get: %f", c.MaxMoney, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "MaxAmount", "int32").(int32); val != c.MaxAmount {
			t.Fatalf("Restore from redis *max amount* fail. Target: %d, Get: %d", c.MaxAmount, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "P", "float32").(float32); val != c.P {
			t.Fatalf("Restore from redis *p* fail. Target: %f, Get: %f", c.P, val)
		}
	}

}
