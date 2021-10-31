package common

import (
	"envelop-rain/db"
	"math/rand"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func SetRandomSeed() {
	rand.Seed(time.Now().Unix())
}

func GetMin(a float32, b float32) float32 {
	if a <= b {
		return a
	}
	return b
}

func Rand() float32 {
	return float32(rand.Intn(101)) / 100
}

func ConvertString(val string, datatype string) interface{} {
	switch datatype {
	case "int":
		ival, _ := strconv.Atoi(val)
		return ival
	case "int32":
		ival, _ := strconv.ParseInt(val, 10, 32)
		return int32(ival)
	case "int64":
		ival, _ := strconv.ParseInt(val, 10, 64)
		return ival
	case "float32":
		fval, _ := strconv.ParseFloat(val, 32)
		return float32(fval)
	case "float64":
		fval, _ := strconv.ParseFloat(val, 64)
		return fval
	case "string":
		return val
	default:
		return val
	}
}

func GetRedPacketsByUID(mysql *gorm.DB, uid int32) ([]*db.RedPacket, error) {
	var packets []*db.RedPacket
	conditions := map[string]interface{}{
		"user_id": uid,
	}
	if err := mysql.Where(conditions).Order("timestamp").Find(&packets).Error; err != nil {
		return nil, err
	}
	return packets, nil
}
