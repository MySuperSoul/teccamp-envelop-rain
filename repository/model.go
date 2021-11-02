/*
 * @Author: your name
 * @Date: 2021-11-02 18:58:08
 * @LastEditTime: 2021-11-02 19:20:51
 * @LastEditors: your name
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/repository/model.go
 */
package repository

type User struct {
	UserID  int32 `gorm:"primaryKey"`
	Amount  int32
	Balance float32
}

type RedPacket struct {
	PacketID  int64 `gorm:"primaryKey"`
	UserID    int32 `gorm:"index:UserID"`
	Value     float32
	Opened    bool
	Timestamp int64
}

func (p *RedPacket) JsonFormat() map[string]interface{} {
	if p.Opened {
		return map[string]interface{}{"envelop_id": p.PacketID, "value": int32(p.Value * 100), "opened": p.Opened, "snatch_time": p.Timestamp}
	}
	return map[string]interface{}{"envelop_id": p.PacketID, "opened": p.Opened, "snatch_time": p.Timestamp}
}
