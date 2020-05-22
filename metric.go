package metric

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Metric struct {
	name      string
	timeStart time.Time
	timeDur   time.Duration
}

type Params struct {
	Sleep     int
	BufferCnt int
	BufferWrn int
	Url       string
	App       string
	Test      bool
}

type Queue struct {
	sync.RWMutex
	values map[string]float64
}

var (
	p = Params{
		Sleep:     5,
		BufferCnt: 100,
		BufferWrn: 50,
		Url:       "127.0.0.1",
		App:       "default",
	}
	q = Queue{
		values: make(map[string]float64, 0),
	}
)

func Init(pp Params) {

	// Init params
	if pp.BufferCnt > 0 {
		p.BufferCnt = pp.BufferCnt
	}
	if pp.BufferWrn > 0 {
		p.BufferWrn = pp.BufferWrn
	}
	if pp.Url != "" {
		p.Url = pp.Url
	}
	if pp.Sleep <= 0 {
		p.Sleep = pp.Sleep
	}
	q = Queue{
		values: make(map[string]float64, pp.BufferCnt),
	}

	// Send queue
	go func() {
		var err error
		for {
			time.Sleep(time.Duration(p.Sleep) * time.Second)
			func() {
				q.Lock()
				defer q.Unlock()
				if len(q.values) > p.BufferWrn {
					log.Printf("cnt message in metric stack is too big %d", len(q.values))
				}
				var s string
				for k, v := range q.values {
					s += fmt.Sprintf("%s_%s %f\n", p.App, k, v)
				}
				err = sendMany(s)
				if err != nil {
					panic(err)
				}
				for k := range q.values {
					delete(q.values, k)
				}
			}()
		}
	}()
}

func Start(name string) Metric {
	m := Metric{
		name:      name,
		timeStart: time.Now(),
	}

	return m
}

func (m *Metric) Stop() {
	m.timeDur = time.Now().Sub(m.timeStart)
	put(m.name+"_seconds", m.timeDur.Seconds())
}

func (m *Metric) Records(value int) {
	put(m.name+"_records", float64(value))
}

func (m *Metric) SubMetric(key string, value float64) {
	put(m.name+"_"+key, value)
}

func put(key string, value float64) {
	q.Lock()
	defer q.Unlock()
	q.values[key] = value
}

func sendMany(plain string) error {
	data := []byte(plain)
	return send(data)
}

func sendOne(key string, value float64) error {
	data := []byte(fmt.Sprintf("%s %f\n", key, value))
	return send(data)
}

func send(data []byte) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/metrics/job/%s", p.Url, p.App), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
