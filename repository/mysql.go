/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-02 19:19:49
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/db/mysql.go
 */
package repository

import (
	"fmt"

	"envelop-rain/configs"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	UserName, PassWord, DB_Name, MysqlAddr string
}

func initDBConfig() DBConfig {
	var db_config DBConfig = DBConfig{
		configs.GlobalConfig.GetString("Database.UserName"),
		configs.GlobalConfig.GetString("Database.Password"),
		configs.GlobalConfig.GetString("Database.DB_Name"),
		configs.GlobalConfig.GetString("Database.MysqlAddr")}

	return db_config
}

func (cfg *DBConfig) formatDBConfig() string {
	var MYSQL_DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.UserName,
		cfg.PassWord,
		cfg.MysqlAddr,
		cfg.DB_Name)
	return MYSQL_DSN
}

func GetMySQLCursor() *gorm.DB {
	db_config := initDBConfig()
	MYSQL_DSN := db_config.formatDBConfig()
	db, err := gorm.Open(mysql.Open(MYSQL_DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Connection to mysql fail")
		panic(err)
	}
	log.Info("MySQL connect success.")
	return db
}

func CloseMySQL(db *gorm.DB) {
	sql, _ := db.DB()
	sql.Close()
}

func GenerateTables(db *gorm.DB) {
	db.AutoMigrate(&DBSysConfig{})
	db.AutoMigrate(&RedPacket{})
}

func CreatePacket(data map[string]interface{}, db *gorm.DB) {
	packet := RedPacket{
		PacketID:  int64(data["packet_id"].(float64)),
		UserID:    int32(data["uid"].(float64)),
		Value:     0,
		Opened:    false,
		Timestamp: int64(data["timestamp"].(float64)),
	}
	db.Create(&packet)
}

func UpdatePacket(data map[string]interface{}, db *gorm.DB) {
	var packet RedPacket
	db.Model(&RedPacket{}).Where("packet_id = ?", int64(data["packet_id"].(float64))).First(&packet)
	packet.Opened = true
	packet.Value = int32(data["value"].(float64))
	db.Save(&packet)
}

func UpdateRemainToDB(data map[string]interface{}, db *gorm.DB) {
	var remain DBSysConfig
	db.Model(&DBSysConfig{}).First(&remain)
	remain.RemainMoney -= int64(data["money"].(float64))
	remain.RemainNum -= int32(data["num"].(float64))
	db.Save(&remain)
}

func SetRemainToDB(config *configs.SystemConfig, db *gorm.DB) {
	remain := DBSysConfig{
		ID:          1,
		RemainNum:   config.TotalNum,
		RemainMoney: config.TotalMoney,
		TotalMoney:  config.TotalMoney,
		TotalNum:    config.TotalNum,
		MaxAmount:   config.MaxAmount,
		P:           config.P,
	}
	db.Create(&remain)
}
