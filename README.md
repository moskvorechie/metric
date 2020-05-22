# Metrics Go client for Prometheus  
Usage
```golang
metric.Init(Params{
	BufferCnt: 100,                // init map buffer count 
	BufferWrn: 50,                 // logs warn if count in buffer more then
	Sleep:     1,                  // send buffer to prom each seconds
	Url:       "http://127.0.0.1", // prom url
	Test:      false,              // for test prevent real send to prom
	App:       "app1",             // prefix
})

// Time metric
x = metric.Start("test1")
time.Sleep(1 * time.Second)
x.Stop()
// Result: app1_test1_seconds: 1.00

// Count
x.Records(100)
// Result: app1_test1_records: 100

// Custom sub metric
x.SubMetric("sub1", 0.1)
// Result: app1_test1_sub1: 0.1
```
