package metric

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Metric struct {
	name   string
	pusher *push.Pusher
}

func New(name string, uri string) Metric {
	return Metric{
		name:   name,
		pusher: push.New(uri, name),
	}
}

func (m *Metric) Pusher() *push.Pusher {
	return m.pusher
}

func (m *Metric) SafePush() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	err := m.pusher.Add()
	if err != nil {
		fmt.Println(err)
	}
}

func (m *Metric) AddGauge(name string) prometheus.Gauge {
	g := prometheus.NewGauge(prometheus.GaugeOpts{Name: name})
	m.pusher.Collector(g)
	return g
}

func (m *Metric) AddSummary(name string) prometheus.Summary {
	g := prometheus.NewSummary(prometheus.SummaryOpts{Name: name})
	m.pusher.Collector(g)
	return g
}

func (m *Metric) AddHistogram(name string) prometheus.Histogram {
	g := prometheus.NewHistogram(prometheus.HistogramOpts{Name: name})
	m.pusher.Collector(g)
	return g
}

func (m *Metric) AddCounter(name string) prometheus.Counter {
	g := prometheus.NewCounter(prometheus.CounterOpts{Name: name})
	m.pusher.Collector(g)
	return g
}
