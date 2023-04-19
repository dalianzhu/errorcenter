`errorcenter` handles errors uniformly in a simple way

User can predefine error types and handling.

When the `Error` function is called, the predefined processing or default processing will be triggered

A sliding counter is provided for the default error processing

## Instructions
```go
var MyErr = ErrType("my error")
var DefaultErr = ErrType("default error")

func init(){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

    // Initialize the whole frame, inactive counters and tickers will be recycled every 15 seconds
	Init(ctx, new(l))

    // Define a handler function for MyErr type
	RegisterHandler(MyErr, func(errType ErrType, msg string, payload map[string]any) {		
        name:= payload["name"] // Make sure there is a name in this error type's payload 

        // Define a key associated with name
        key := fmt.Sprintf("errType:%v:%v", errType, name)

        // Get a Counter object, this object will be cached
        counter := GetCounter(key)
        counter.Add(1)

        times := counter.Sum()
        // get a tickerDo object, periodically print logs
        GetTicker(key, 1).Do(func() {
	    	log.Printf("%v, %v, %v, %v\n", MyErr, msg, payload, times)
		})
	})
}

func run(){
    err:= doSth()
    if err!=nil{
	    Error(DefaultErr, fmt.Sprintf("error:%v", err), nil)
        /*
        output:
        errType:default error,msg:error:ERRORINFO,payload:map[]

        or if frequently called

        errType:default error,msg:error:ERRORINFO,payload:map[],counts in the past 5 mins:27283
        */
    }	

    // or trigger predefined error types
	Error(MyErr, "my error", map[string]any{
		"name": "pig",
	})
    /*
    output:
    my error, my error msg, map[name:pig], 15
    */
}
	



```