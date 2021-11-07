/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-02 20:59:03
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/common/util.go
 */
package common

import (
	"math/rand"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func SetRandomSeed() {
	rand.Seed(time.Now().Unix())
}

func GetMin(a int32, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}

func Rand() float32 {
	return float32(rand.Intn(101)) / 100
}

func ConvertString(val string, datatype string) interface{} {
	if val == "" {
		log.Error("(func ConvertString) Null String from Post request!")
	}

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
	case "bool":
		bval, _ := strconv.ParseBool(val)
		return bval
	case "string":
		return val
	default:
		return val
	}
}

func GetRedPacket(remain_num int32, remain_money int64, min_money int32, max_money int32) int32 {
	if remain_num == 1 {
		return GetMin(int32(remain_money), max_money)
	}
	mean_money := int32(remain_money / int64(remain_num))
	max_money = GetMin(max_money, 2*mean_money-min_money)
	money := min_money + rand.Int31n(max_money-min_money+1)
	return money
}
