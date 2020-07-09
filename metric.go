package metric

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	BufferWrn int
	Url       string
	App       string
	Test      bool
	isInit    bool
}

type rec struct {
	name  string
	value float64
}

type queue struct {
	sync.RWMutex
	values []rec
}

var (
	p = Params{
		Sleep:     5,
		BufferWrn: 50,
		Url:       "127.0.0.1",
		App:       "default",
	}
	q = queue{
		values: make([]rec, 0),
	}
)

func Init(pp Params) {

	// Init params
	if pp.BufferWrn > 0 {
		p.BufferWrn = pp.BufferWrn
	}
	if pp.Url != "" {
		p.Url = pp.Url
	}
	if pp.Sleep <= 0 {
		p.Sleep = pp.Sleep
	}
	if pp.App != "default" {
		p.App = pp.App
	}
	q = queue{
		values: make([]rec, 0),
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
					log.Printf("cnt message in Metric stack is too big %d", len(q.values))
				}
				var s string
				for _, r := range q.values {
					s += fmt.Sprintf("%s_%s %f\n", p.App, r.name, r.value)
				}
				err = sendMany(s)
				if err != nil {
					panic(err)
				}
				q.values = nil
			}()
		}
	}()

	p.isInit = true
}

func Start(name string) Metric {
	if !p.isInit {
		panic("Metric is not init")
	}
	m := Metric{
		name:      name,
		timeStart: time.Now(),
	}

	return m
}

func (m *Metric) Stop() {
	m.timeDur = time.Now().Sub(m.timeStart)
	put(rec{
		name:  m.name + "_seconds",
		value: m.timeDur.Seconds(),
	})
}

func (m *Metric) Records(value interface{}) {
	var res float64
	switch value.(type) {
	case int:
		res = float64(value.(int))
	case int64:
		res = float64(value.(int64))
	case string:
		i, err := strconv.Atoi(value.(string))
		if err != nil {
			panic(err)
		}
		res = float64(i)
	}

	put(rec{
		name:  m.name + "_records",
		value: res,
	})
}

func (m *Metric) SubMetric(key string, value float64) {
	put(rec{
		name:  m.name + "_" + key,
		value: value,
	})
}

func put(r rec) {
	q.Lock()
	defer q.Unlock()
	q.values = append(q.values, r)
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
