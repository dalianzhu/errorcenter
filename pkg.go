package errorcenter

import (
	"context"
	"sync"
	"time"
)

func Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 15):
			globalCounterMap.Range(func(key, value any) bool {
				v := value.(*SlidingWindowCounter)
				if v.Sum() <= 0 {
					globalCounterMapLk.Lock()
					globalCounterMap.Delete(key)
					globalCounterMapLk.Unlock()
				}
				return true
			})
			globalTickerMap.Range(func(key, value any) bool {
				v := value.(*TickerDo)
				if time.Now().Unix()-v.LastCallTime() > 300 {
					globalTickerLk.Lock()
					globalTickerMap.Delete(key)
					globalTickerLk.Unlock()
				}
				return true
			})
		}
	}
}

var globalCounterMap sync.Map
var globalCounterMapLk sync.Mutex

// GetCounter if the key does not exist, create a counter.
// The counter will be checked periodically, and the inactive ones will be recycled.
func GetCounter(name string) *SlidingWindowCounter {
	v, ok := globalCounterMap.Load(name)
	if ok {
		return v.(*SlidingWindowCounter)
	}

	globalCounterMapLk.Lock()
	defer globalCounterMapLk.Unlock()

	v, ok = globalCounterMap.Load(name)
	if ok {
		return v.(*SlidingWindowCounter)
	}
	c := NewSlidingWindowCounter(5)
	globalCounterMap.Store(name, c)
	return c
}

var globalTickerMap sync.Map
var globalTickerLk sync.Mutex

// GetTicker if the key does not exist, create a TickerDo obj.
// The item will be checked periodically, and the inactive ones will be recycled.
func GetTicker(name string, duSecs int) *TickerDo {
	v, ok := globalTickerMap.Load(name)
	if ok {
		return v.(*TickerDo)
	}
	globalTickerLk.Lock()
	defer globalTickerLk.Unlock()

	v, ok = globalTickerMap.Load(name)
	if ok {
		return v.(*TickerDo)
	}
	tk := NewTickDo(duSecs)
	globalTickerMap.Store(name, tk)
	return tk
}
