package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type SystemConfig struct {
	TotalMoney, MinMoney, MaxMoney, P float32
	TotalNum, MaxAmount               int32
}

/*
params:
	file_path: The path of config json file
return:
	SystemConfig
*/
func GenerateConfigFromFile(file_path string) SystemConfig {
	json_file, err := os.Open(file_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer json_file.Close()

	bytevalues, _ := ioutil.ReadAll(json_file)
	var config SystemConfig
	json.Unmarshal(bytevalues, &config)

	return config
}

func SetConfigToRedis(config *SystemConfig, redisdb *redis.Client) {
	err := redisdb.MSet(
		"TotalMoney", config.TotalMoney,
		"TotalNum", config.TotalNum,
		"RemainMoney", config.TotalMoney,
		"RemainNum", config.TotalNum,
		"MinMoney", config.MinMoney,
		"MaxMoney", config.MaxMoney,
		"MaxAmount", config.MaxAmount,
		"P", config.P,
	).Err()
	if err != nil {
		log.Fatal("Set config to redis error.")
		panic(err)
	}
}
