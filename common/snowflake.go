package common

import (
	"errors"
	"math"
	"sync"
	"time"
)

const (
	workerIDBits     = uint64(2) // workerID
	dataCenterIDBits = uint64(2) // dataCenterID
	sequenceBits     = uint64(12)

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits)
	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft     = uint8(22)
	dataLeft     = uint8(17)
	workLeft     = uint8(12)
	epoch        = int64(1589923200000)
	MaxQueueSize = 50000
)

type Worker struct {
	mu           sync.Mutex
	LastStamp    int64
	WorkerID     int64
	DataCenterID int64
	Sequence     int64
}

type IDProducer struct {
	IDs      chan (int64)
	Producer *Worker
	TotalNum int
}

func NewProducer(worker_id, datacenter_id int64, total_num int) *IDProducer {
	return &IDProducer{
		IDs:      make(chan int64, MaxQueueSize),
		TotalNum: total_num,
		Producer: &Worker{
			LastStamp:    0,
			WorkerID:     worker_id,
			DataCenterID: datacenter_id,
			Sequence:     0,
		},
	}
}

func (p *IDProducer) StartProducePacketID() {
	for {
		id, err := p.Producer.NextID()
		if err != nil {
			continue
		}
		p.IDs <- id
	}
}

func (w *Worker) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (w *Worker) NextID() (int64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.nextID()
}

func (w *Worker) nextID() (int64, error) {
	timeStamp := w.getMilliSeconds()
	if timeStamp < w.LastStamp {
		return 0, errors.New("time is moving backwards,waiting until")
	}

	if w.LastStamp == timeStamp {

		w.Sequence = (w.Sequence + 1) & maxSequence

		if w.Sequence == 0 {
			for timeStamp <= w.LastStamp {
				timeStamp = w.getMilliSeconds()
			}
		}
	} else {
		w.Sequence = 0
	}

	w.LastStamp = timeStamp
	id := ((timeStamp - epoch) << timeLeft) | (w.DataCenterID << dataLeft) | (w.WorkerID << workLeft) | w.Sequence
	return int64(id) % int64(math.Pow10(15)), nil
}
