package metric

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Metric struct {
	stated    bool
	name      string
	timeStart time.Time
	timeDur   time.Duration
}

type Params struct {
	Sleep      int
	BufferWrn  int
	Url        string
	App        string
	Instance   string
	Test       bool
	Debug      bool
	HttpClient *http.Client
	isInit     bool
}

type RecType string

const (
	RecTypeHistogram = "histogram"
	RecTypeCounter   = "counter"
	RecTypeGauge     = "gauge"
	RecTypeSummary   = "summary"
)

type rec struct {
	Name  string
	Value float64
	Type  string
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
	if pp.Instance != "default" {
		p.Instance = pp.Instance
	}
	p.Test = pp.Test
	p.Debug = pp.Debug
	q = queue{
		values: make([]rec, 0),
	}

	p.HttpClient = &http.Client{}

	// Send queue
	go func() {
		var err error
		for {
			time.Sleep(time.Duration(p.Sleep) * time.Second)
			func() {
				q.Lock()
				defer q.Unlock()
				if len(q.values) <= 0 {
					return
				}
				if p.Debug && len(q.values) > p.BufferWrn {
					log.Printf("cnt message in Metric stack is too big %d", len(q.values))
				}
				t := make(map[string]bool, 0)
				var s string
				for _, r := range q.values {
					// Set type one time
					if b, ok := t[fmt.Sprintf("%s_%s", p.App, r.Name)]; !ok || b == false {
						s += fmt.Sprintf("# TYPE %s_%s %s\n", p.App, r.Name, r.Type)
						t[fmt.Sprintf("%s_%s", p.App, r.Name)] = true
					}
					// Set value
					s += fmt.Sprintf("%s_%s %f\n", p.App, r.Name, r.Value)
				}
				if len(s) <= 0 {
					return
				}
				err = sendStr(s)
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
		stated:    true,
		name:      name,
		timeStart: time.Now(),
	}

	return m
}

func (m *Metric) Stop() {
	if !m.stated {
		return
	}
	m.timeDur = time.Now().Sub(m.timeStart)
	put(rec{
		Name:  m.name + "_seconds",
		Value: m.timeDur.Seconds(),
		Type:  RecTypeGauge,
	})
}

// Save records count
func (m *Metric) Records(value interface{}) {
	if !m.stated {
		return
	}
	put(rec{
		Name:  m.name + "_records",
		Value: toFloat(value),
		Type:  RecTypeGauge,
	})
}

// Save records count
func Values() []rec {
	return q.values
}

// Additional metric with same app and Name
func (m *Metric) SubMetric(key string, value interface{}) {
	if !m.stated {
		return
	}
	put(rec{
		Name:  m.name + "_" + key,
		Value: toFloat(value),
		Type:  RecTypeGauge,
	})
}

// Additional metric with same app and Name
func (m *Metric) Custom(key string, value interface{}, sType string) {
	if !m.stated {
		return
	}
	switch sType {
	case RecTypeHistogram:
	case RecTypeGauge:
	case RecTypeSummary:
	case RecTypeCounter:
	default:
		panic("no type found")
	}
	put(rec{
		Name:  m.name + "_" + key,
		Value: toFloat(value),
		Type:  sType,
	})
}

func put(r rec) {
	q.Lock()
	defer q.Unlock()
	q.values = append(q.values, r)
}

func sendStr(plain string) error {
	data := []byte(plain)
	return send(data)
}

func send(data []byte) error {
	if p.Test {
		return nil
	}
	if len(data) <= 0 {
		return nil
	}

	uri := fmt.Sprintf("%s/metrics/job/%s", p.Url, p.App)
	if p.Instance != "" {
		uri = fmt.Sprintf("%s/metrics/job/%s/instance/%s", p.Url, p.App, p.Instance)
	}
	if p.Debug {
		fmt.Println(uri, strings.TrimSpace(string(data)))
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := p.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if p.Debug {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println("code:", resp.StatusCode, "body:", string(body))
	}
	return nil
}

func toFloat(value interface{}) float64 {
	var res float64
	switch value.(type) {
	case int:
		res = float64(value.(int))
	case int32:
		res = float64(value.(int32))
	case int64:
		res = float64(value.(int64))
	case float32:
		res = float64(value.(float32))
	case float64:
		res = value.(float64)
	case string:
		i, _ := strconv.Atoi(value.(string))
		res = float64(i)
	default:
		panic("This type not supported")
	}

	return res
}
