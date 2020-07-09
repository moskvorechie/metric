package metric

import (
	"testing"
	"time"
)

func init() {
	Init(Params{
		BufferWrn: 50,
		Sleep:     1,
		Url:       "http://127.0.0.1",
		Test:      true,
		App:       "test",
	})
}

func TestMetric_Stop(t *testing.T) {
	x := Start("test2")
	time.Sleep(100 * time.Millisecond)
	x.Stop()
	if len(q.values) <= 0 {
		t.Fatal()
	}
	x.Records(100)
	if len(q.values) <= 1 {
		t.Fatal()
	}
	if v := q.values[1]; v.value != 100 {
		t.Fatal()
	}
	x.SubMetric("s1", 0.1)
	if len(q.values) <= 2 {
		t.Fatal()
	}
	if v := q.values[2]; v.value != 0.1 {
		t.Fatal()
	}
	time.Sleep(15 * time.Second)
}
