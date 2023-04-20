package errorcenter

import (
	"context"
	"fmt"
)

// Init pass in the log object, start the error center
func Init(ctx context.Context, globalLog interface {
	Errorf(string, ...any)
}) {
	GlobalLog = globalLog
	go Start(ctx)
}

// Error trigger a specific type of error and hand it over to the error center
func Error(errType ErrType, msg string, payload map[string]any) {
	customHandler, ok := codeHandlerMap[errType]
	if ok {
		customHandler(errType, msg, payload)
	} else {
		DefaultHandler(errType, msg, payload)
	}
}

// RegisterHandler register a specific type of error handler function
func RegisterHandler(errType ErrType, fn CodeHandler) {
	codeHandlerMap[errType] = fn
}

// ErrType represents each different type of error in this system
type ErrType string

var GlobalLog interface {
	Errorf(string, ...any)
}

var codeHandlerMap = make(map[ErrType]CodeHandler)

// CodeHandler users can customize the handling of certain types of errors
type CodeHandler func(errType ErrType, msg string, payload map[string]any)

// DefaultHandler defines general processing,
// it will record the number of errors of the same type within 5 minutes and print messages regularly
var DefaultHandler = func(errType ErrType, msg string, payload map[string]any) {
	key := fmt.Sprintf("errType:%v", errType)
	counter := GetCounter(key)
	counter.Add(1)
	times := counter.Sum()
	if times == 1 {
		GlobalLog.Errorf("errType:%v,msg:%v,payload:%v",
			errType, msg, payload)
	} else {
		GetTicker(key, 1).Do(func() {
			GlobalLog.Errorf("errType:%v,msg:%v,payload:%v,counts in the past 5 mins:%v",
				errType, msg, payload, times)
		})
	}
}
