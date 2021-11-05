/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-01 16:01:17
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/configs/config.go
 */
package configs

import (
	"path"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var GlobalConfig *viper.Viper

func init() {
	log.Debug("Loading configuration...")
	GlobalConfig = initConfig()
	go watchConfig()
}

func initConfig() *viper.Viper {
	var abPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		abPath = path.Dir(filename)
	}
	log.Info(abPath, filename)
	Config := viper.New()
	Config.SetConfigType("yaml")
	Config.SetConfigName("config")
	Config.AddConfigPath(abPath)
	Config.AddConfigPath(".")
	Config.AddConfigPath("./configs")
	Config.AutomaticEnv()
	Config.SetEnvPrefix("Envelop_Rain")
	replacer := strings.NewReplacer(".", "_")
	Config.SetEnvKeyReplacer(replacer)
	if err := Config.ReadInConfig(); err != nil {
		log.Error(err)
		log.Fatal("Faild to read the configuration.")
	}
	return Config
}

// 监控配置文件变化并热加载程序 TO_DO
func watchConfig() {
	GlobalConfig.WatchConfig()
	GlobalConfig.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}

type SystemConfig struct {
	TotalMoney                              int64
	P                                       float32
	MinMoney, MaxMoney, TotalNum, MaxAmount int32
}

/*
return:
	SystemConfig
*/
func GenerateConfigFromFile() SystemConfig {
	config := SystemConfig{
		GlobalConfig.GetInt64("common.TotalMoney"),
		float32(GlobalConfig.GetFloat64("common.P")),
		GlobalConfig.GetInt32("common.MinMoney"),
		GlobalConfig.GetInt32("common.MaxMoney"),
		GlobalConfig.GetInt32("common.TotalNum"),
		GlobalConfig.GetInt32("common.MaxAmount")}
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
