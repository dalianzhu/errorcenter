package errorcenter

import (
	"context"
	"fmt"
	"log"
	"testing"
)

type l struct{}

func (l *l) Errorf(fmt string, args ...any) {
	log.Printf(fmt, args...)
}
func TestError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Init(ctx, new(l))
	var MyErr = ErrType("my error")
	var DefaultErr = ErrType("default error")

	RegisterHandler(MyErr, func(errType ErrType, msg string, payload map[string]any) {
		log.Printf("%v, %v, %v\n", MyErr, msg, payload)
	})

	Error(MyErr, "my error msg", map[string]any{
		"name": "pig",
	})
	for i := 0; i < 10; i++ {
		Error(DefaultErr, fmt.Sprintf("##default error:%v##", i), nil)
	}
}
