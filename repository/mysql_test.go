package repository

import (
	"testing"
)

func TestMysql(t *testing.T) {
	db := GetMySQLCursor()

	if db == nil {
		t.Failed()
	}

	GenerateTables(db)
	// user表中插入一条记录
	packet := RedPacket{PacketID: 12345, UserID: 1234, Value: 30, Opened: false, Timestamp: 123456}
	db.Create(&packet)
	// user表查找
	var p RedPacket
	db.Where(&packet).First(&p)
	if p != packet {
		t.Fatal("select from user failed")
	}
	db.Delete(&packet)

	CloseMySQL(db)
}
