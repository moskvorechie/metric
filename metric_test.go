package metric_test

import (
	"github.com/moskvorechie/metric"
	"testing"
	"time"
)

func init() {
	metric.Init(metric.Params{
		BufferWrn: 50,
		Sleep:     1,
		Url:       "http://127.0.0.1",
		Test:      true,
		App:       "test",
	})
}

func TestMetric_Stop(t *testing.T) {
	x := metric.Start("test2")
	time.Sleep(100 * time.Millisecond)
	x.Stop()
	if len(metric.Values()) <= 0 {
		t.Fatal()
	}
	x.Records(100)
	if len(metric.Values()) <= 1 {
		t.Fatal()
	}
	if v := metric.Values()[1]; v.Value != 100 {
		t.Fatal()
	}
	x.SubMetric("s1", 0.1)
	if len(metric.Values()) <= 2 {
		t.Fatal()
	}
	if v := metric.Values()[2]; v.Value != 0.1 {
		t.Fatal()
	}
	time.Sleep(15 * time.Second)
}
