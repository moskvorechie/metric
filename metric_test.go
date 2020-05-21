package metric

import (
	"testing"
	"time"
)

func init() {
	Init(Params{
		BufferCnt: 100,
		BufferWrn: 50,
		Sleep:     1,
		Url:       "http://127.0.0.1",
		Test:      false,
		App:       "test",
	})
}

func TestMetric_Stop(t *testing.T) {
	x := Start("test2")
	time.Sleep(100 * time.Millisecond)
	x.Stop()
	if len(q.values) <= 0 {
		t.Fail()
	}
	x.Records(100)
	if len(q.values) <= 1 {
		t.Fail()
	}
	if v, ok := q.values["test_records"]; !ok || v != 100 {
		t.Fail()
	}
	x.SubMetric("s1", 0.1)
	if len(q.values) <= 2 {
		t.Fail()
	}
	time.Sleep(15 * time.Second)
}
