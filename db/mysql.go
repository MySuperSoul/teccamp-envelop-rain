package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	USERNAME   = "root"
	PASSWORD   = "yourpassword" // use your password here
	DB_NAME    = "envelop-rain"
	MYSQL_ADDR = "127.0.0.1:3306"
)

var MYSQL_DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USERNAME, PASSWORD, MYSQL_ADDR, DB_NAME)

func GetMySQLCursor() *gorm.DB {
	db, err := gorm.Open(mysql.Open(MYSQL_DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Connection to mysql fail")
		panic(err)
	}
	log.Info("MySQL connect success.")
	return db
}

func GenerateTables(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&RedPacket{})
}
