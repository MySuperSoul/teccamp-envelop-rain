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
	user := User{UserID: 111111, Amount: 0, Balance: 0.}
	db.Create(&user)
	// user表查找
	var userDB User
	db.Where(&user).First(&userDB)
	if userDB != user {
		t.Fatal("select from user failed")
	}
	db.Delete(&user)

	CloseMySQL(db)
}
