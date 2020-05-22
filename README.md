# Metrics Go client for Prometheus  
Usage
```golang
metric.Init(Params{
    BufferCnt: 100,
    BufferWrn: 50,
    Sleep:     1,
    Url:       "http://127.0.0.1",
    Test:      false,
    App:       "test",
})

// Time metric
x = metric.Start("test1")
time.Sleep(1 * time.Second)
x.Stop()
// Result: test1_seconds: 1.00

// Count
x.Records(100)
// Result: test1_records: 100

// Custom sub metric
x.SubMetric("sub1", 0.1)
// Result: test1_sub1: 0.1
```