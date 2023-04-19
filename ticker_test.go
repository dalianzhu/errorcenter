package errorcenter

import (
	"log"
	"testing"
	"time"
)

func TestTickerDo_Do(t *testing.T) {
	tk := NewTickDo(1)
	for {
		go tk.Do(func() {
			log.Printf("time:%v\n", time.Now())
		})
	}
}
