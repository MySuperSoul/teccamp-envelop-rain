package configs

import (
	"encoding/json"
	"envelop-rain/db"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestGenerateConfigFromFile(t *testing.T) {
	filename := "tmp.json"
	cases := []SystemConfig{
		{1., 1., 1., 1., 1, 1},
		{1., 10., 100., 1000., 1234, 4},
		{3.4, 10.5, 34.5, 12.3, 9, 8},
	}

	for _, c := range cases {
		// create file and write json content into it
		fileptr, _ := os.Create(filename)
		encoder := json.NewEncoder(fileptr)
		err := encoder.Encode(c)
		if err != nil {
			fmt.Println(err)
		}
		fileptr.Close()

		// Generate from config method
		config := GenerateConfigFromFile(filename)
		if c != config {
			os.Remove(filename)
			t.Fatalf("Generate config failed")
		}

		// finally delete the temp file
		os.Remove(filename)
	}
}

func TestSetConfigToRedis(t *testing.T) {
	redisdb := db.GetRedisClient()
	defer redisdb.Close()

	cases := []SystemConfig{
		{1., 1., 1., 1., 1, 1},
		{1., 10., 100., 1000., 1234, 4},
		{3.4, 10.5, 34.5, 12.3, 9, 8},
	}

	for _, c := range cases {
		// first set to redis
		SetConfigToRedis(&c, redisdb)
		vals, err := redisdb.MGet(
			"TotalMoney", "TotalNum", "RemainMoney",
			"RemainNum", "MinMoney", "MaxMoney", "MaxAmount", "P").Result()
		if err != nil {
			panic(err)
		}

		if val, _ := strconv.ParseFloat(vals[0].(string), 32); float32(val) != c.TotalMoney {
			t.Fatalf("Restore from redis *total money* fail. Target: %f, Get: %f", c.TotalMoney, float32(val))
		}
		if val, _ := strconv.Atoi(vals[1].(string)); val != int(c.TotalNum) {
			t.Fatalf("Restore from redis *total num* fail. Target: %d, Get: %d", c.TotalNum, val)
		}
		if val, _ := strconv.ParseFloat(vals[2].(string), 32); float32(val) != c.TotalMoney {
			t.Fatalf("Restore from redis *remain money* fail. Target: %f, Get: %f", c.TotalMoney, float32(val))
		}
		if val, _ := strconv.Atoi(vals[3].(string)); val != int(c.TotalNum) {
			t.Fatalf("Restore from redis *remain num* fail. Target: %d, Get: %d", c.TotalNum, val)
		}
		if val, _ := strconv.ParseFloat(vals[4].(string), 32); float32(val) != c.MinMoney {
			t.Fatalf("Restore from redis *min money* fail. Target: %f, Get: %f", c.MinMoney, float32(val))
		}
		if val, _ := strconv.ParseFloat(vals[5].(string), 32); float32(val) != c.MaxMoney {
			t.Fatalf("Restore from redis *max money* fail. Target: %f, Get: %f", c.MaxMoney, float32(val))
		}
		if val, _ := strconv.Atoi(vals[6].(string)); val != int(c.MaxAmount) {
			t.Fatalf("Restore from redis *max amount* fail. Target: %d, Get: %d", c.MaxAmount, val)
		}
		if val, _ := strconv.ParseFloat(vals[7].(string), 32); float32(val) != c.P {
			t.Fatalf("Restore from redis *p* fail. Target: %f, Get: %f", c.P, float32(val))
		}
	}

}
