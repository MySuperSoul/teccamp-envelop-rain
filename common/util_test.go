package common

import (
	"testing"
)

func TestRand(t *testing.T) {
	loop_times := 100
	for i := 0; i < loop_times; i++ {
		ratio := Rand()
		if ratio < 0 || ratio > 1 {
			t.Fatalf("Out of boundary, value is %f", ratio)
		}
	}
}
