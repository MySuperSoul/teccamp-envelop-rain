/*
 * @Author: your name
 * @Date: 2021-11-01 14:56:45
 * @LastEditTime: 2021-11-02 19:19:27
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/db/db_test.go
 */
package repository

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
		{TotalMoney: 1, MinMoney: 1, MaxMoney: 1, P: 1., TotalNum: 1, MaxAmount: 1},
		{TotalMoney: 1, MinMoney: 10, MaxMoney: 100, P: 1000., TotalNum: 1234, MaxAmount: 4},
		{TotalMoney: 34, MinMoney: 105, MaxMoney: 345, P: 12.3, TotalNum: 9, MaxAmount: 8}}

	for _, c := range cases {
		// first set to redis
		configs.SetConfigToRedis(&c, redisdb)

		if val := GetSingleValueFromRedis(redisdb, "TotalMoney", "int64").(int64); val != c.TotalMoney {
			t.Fatalf("Restore from redis *total money* fail. Target: %d, Get: %d", c.TotalMoney, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "TotalNum", "int32").(int32); val != c.TotalNum {
			t.Fatalf("Restore from redis *total num* fail. Target: %d, Get: %d", c.TotalNum, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "RemainMoney", "int64").(int64); val != c.TotalMoney {
			t.Fatalf("Restore from redis *remain money* fail. Target: %d, Get: %d", c.TotalMoney, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "RemainNum", "int32").(int32); val != c.TotalNum {
			t.Fatalf("Restore from redis *remain num* fail. Target: %d, Get: %d", c.TotalNum, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "MinMoney", "int32").(int32); val != c.MinMoney {
			t.Fatalf("Restore from redis *min money* fail. Target: %d, Get: %d", c.MinMoney, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "MaxMoney", "int32").(int32); val != c.MaxMoney {
			t.Fatalf("Restore from redis *max money* fail. Target: %d, Get: %d", c.MaxMoney, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "MaxAmount", "int32").(int32); val != c.MaxAmount {
			t.Fatalf("Restore from redis *max amount* fail. Target: %d, Get: %d", c.MaxAmount, val)
		}
		if val := GetSingleValueFromRedis(redisdb, "P", "float32").(float32); val != c.P {
			t.Fatalf("Restore from redis *p* fail. Target: %f, Get: %f", c.P, val)
		}
	}

}

func TestRecoverFromRedis(t *testing.T) {
	c := map[string]interface{}{
		"amount":  3,
		"balance": 69.8,
	}
	redisdb := GetRedisClient()
	defer redisdb.Close()
	err := redisdb.HMSet("123", c).Err()
	if err != nil {
		panic(err)
	}

	redisdb.Del("123-wallet")
	redisdb.LPush("123-wallet", 12345, 123456.7)

	if l, _ := redisdb.LLen("123-wallet").Result(); l != 2 {
		t.Fatal("Get length wrong")
	}
}

func TestMysql(t *testing.T) {
	db := GetMySQLCursor()

	if db == nil {
		t.Failed()
	}
	sql, _ := db.DB()
	defer sql.Close()

	GenerateTables(db)
	// user表中插入一条记录
	user := User{UserID: 111111, Amount: 0, Balance: 0.}
	db.Create(&user)
	// user表查找
	var userDB User
	db.Where(&user).First(&userDB)
	if userDB != user {
		t.Fatal("select from user failed")
	}
	db.Delete(&user)
}
