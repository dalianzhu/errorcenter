package errorcenter

import (
	"sync"
	"time"
)

// NewSlidingWindowCounter create a sliding window counter, slides forward every minute
func NewSlidingWindowCounter(len int) *SlidingWindowCounter {
	s := new(SlidingWindowCounter)
	s.counts = make([]int, len)
	return s
}

// SlidingWindowCounter counting the data of the last N minutes by sliding every minute
type SlidingWindowCounter struct {
	counts []int
	sync.RWMutex
	lastUpdate int64
	offset     int
}

// Add add count to the current minute value
func (e *SlidingWindowCounter) Add(count int) {
	e.Lock()
	defer e.Unlock()
	e.counts[e.offset] += count
	e.lastUpdate = time.Now().Unix()
}

// Tick move the pointer to the earliest minute and set the data to 0
//
//	| offset       | offset
//
// [13,4,1] -> [13,0,1]
func (e *SlidingWindowCounter) Tick() {
	e.Lock()
	defer e.Unlock()
	e.offset += 1
	if e.offset >= len(e.counts) {
		e.offset = 0
	}
	e.counts[e.offset] = 0
}

// Sum returns the sum of all values
func (e *SlidingWindowCounter) Sum() int {
	e.Lock()
	defer e.Unlock()
	all := 0
	for _, v := range e.counts {
		all += v
	}
	return all
}

// GetData returns a copy of the original data
func (e *SlidingWindowCounter) GetData() []int {
	e.Lock()
	defer e.Unlock()

	ret := make([]int, len(e.counts))
	for i := 0; i < len(e.counts); i++ {
		offset := e.offset - i
		if offset < 0 {
			offset = len(e.counts) + offset
		}
		ret[len(e.counts)-i-1] = e.counts[offset]
	}
	return ret
}
