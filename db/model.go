package db

type User struct {
	UserID  int32 `gorm:"primaryKey"`
	Amount  int32
	Balance float32
}

type RedPacket struct {
	PacketID  int64 `gorm:"primaryKey"`
	UserID    int32
	Value     float32
	Opened    bool
	Timestamp int64
}
