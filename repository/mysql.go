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

func GetRedPacketsByUID(mysql *gorm.DB, uid int32) ([]*RedPacket, error) {
	var packets []*RedPacket
	conditions := map[string]interface{}{
		"user_id": uid,
	}
	if err := mysql.Where(conditions).Order("timestamp").Find(&packets).Error; err != nil {
		return nil, err
	}
	return packets, nil
}

func GenerateTables(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&RedPacket{})
}
