/*
 * @Author: your name
 * @Date: 2021-11-02 18:58:08
 * @LastEditTime: 2021-11-02 19:20:51
 * @LastEditors: your name
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/repository/model.go
 */
package repository

type DBSysConfig struct {
	ID          int32 `gorm:"primaryKey"`
	RemainNum   int32
	RemainMoney int64
	TotalMoney  int64
	TotalNum    int32
	MaxAmount   int32
	P           float32
}

type RedPacket struct {
	PacketID  int64 `gorm:"primaryKey"`
	UserID    int32 `gorm:"index:UserID"`
	Value     int32
	Opened    bool
	Timestamp int64
}

func (p *RedPacket) ToRedisFormat() map[string]interface{} {
	return map[string]interface{}{
		"userid":    p.UserID,
		"value":     p.Value,
		"opened":    p.Opened,
		"timestamp": p.Timestamp,
	}
}

func (p *RedPacket) JsonFormat() map[string]interface{} {
	if p.Opened {
		return map[string]interface{}{"envelop_id": p.PacketID, "value": p.Value, "opened": p.Opened, "snatch_time": p.Timestamp}
	}
	return map[string]interface{}{"envelop_id": p.PacketID, "opened": p.Opened, "snatch_time": p.Timestamp}
}
