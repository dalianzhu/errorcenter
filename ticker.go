package errorcenter

import (
	"sync/atomic"
	"time"
)

func NewTickDo(duSecs int) *TickerDo {
	t := new(TickerDo)
	t.du = int64(duSecs)
	return t
}

// TickerDo guarantees the minimum interval for the execution of the Do function
type TickerDo struct {
	lastCallTime int64
	du           int64
}

// LastCallTime returns the timestamp of the most recent call
func (t *TickerDo) LastCallTime() int64 {
	return atomic.LoadInt64(&t.lastCallTime)
}

// Do can be called concurrently, and will only be executed once within the duration
// useful in scenarios such as printing a large number of logs
func (t *TickerDo) Do(fn func()) {
	now := time.Now().Unix()
	lastCallTime := atomic.LoadInt64(&t.lastCallTime)
	if now-lastCallTime >= t.du {
		ok := atomic.CompareAndSwapInt64(&t.lastCallTime, lastCallTime, now)
		if ok {
			fn()
		}
		return
	} else {
		return
	}
}
