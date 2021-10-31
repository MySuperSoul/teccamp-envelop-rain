package common

import (
	"math/rand"
	"strconv"
	"time"
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