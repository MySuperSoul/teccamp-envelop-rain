package common

import (
	"sync"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	var wg sync.WaitGroup
	worker := &Worker{
		LastStamp:    0,
		WorkerID:     0,
		DataCenterID: 0,
		Sequence:     0,
	}
	size := 1000000
	ch := make(chan (int64), size)

	wg.Add(size)
	defer close(ch)

	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			id, _ := worker.NextID()
			ch <- id
		}()
	}
	wg.Wait()

	m := make(map[int64]int)
	for i := 0; i < size; i++ {
		id := <-ch
		_, ok := m[id]
		if ok {
			t.Fatalf("repeat id %d\n", id)
			return
		}
		m[id] = i
	}
}
