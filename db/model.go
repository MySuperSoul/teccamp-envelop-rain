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

func (p *RedPacket) JsonFormat() map[string]interface{} {
	if p.Opened {
		return map[string]interface{}{"envelop_id": p.PacketID, "value": p.Value, "opened": p.Opened, "snatch_time": p.Timestamp}
	}
	return map[string]interface{}{"envelop_id": p.PacketID, "opened": p.Opened, "snatch_time": p.Timestamp}
}
