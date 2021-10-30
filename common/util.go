package common

import (
	"math/rand"
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
