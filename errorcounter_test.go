package errorcenter

import (
	"log"
	"testing"
	"time"
)

func TestSlidingWindowCounter_Tick(t *testing.T) {
	s := NewSlidingWindowCounter(5)
	go func() {
		for {
			time.Sleep(time.Minute)
			s.Tick()
		}
	}()
	s.Add(1)
	for {
		log.Printf("all:%v l:%v \n", s.Sum(), s.GetData())
		time.Sleep(time.Second * 15)
	}
}
