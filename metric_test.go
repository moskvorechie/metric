package metric_test

import (
	"github.com/moskvorechie/metric/v2"
	"testing"
)

func TestMetric_New(t *testing.T) {

	m := metric.New("app_test", "http://127.0.0.1:30091")
	defer m.SafePush()

	records := m.AddGauge("db_backup_records_processed")
	records.Set(101.00)

	duration := m.AddGauge("db_backup_duration_seconds")
	duration.Set(5.00)

	//g := m.StartGauge("metric_1")
	//time.Sleep(100 * time.Millisecond)

	//g.StopAndPush()

	//registry := prometheus.NewRegistry()
	//registry.MustRegister(completionTime, duration, records)

	//func init() {
	//	//metric.Init(metric.Params{
	//	//	BufferWrn: 50,
	//	//	Sleep:     1,
	//	//	Url:       "http://127.0.0.1",
	//	//	Test:      true,
	//	//	App:       "test",
	//	//})
	//}

	//x := metric.Start("test2")
	//time.Sleep(100 * time.Millisecond)
	//x.Stop()
	//if len(metric.Values()) <= 0 {
	//	t.Fatal()
	//}
	//x.Records(100)
	//if len(metric.Values()) <= 1 {
	//	t.Fatal()
	//}
	//if v := metric.Values()[1]; v.Value != 100 {
	//	t.Fatal()
	//}
	//x.SubMetric("s1", 0.1)
	//if len(metric.Values()) <= 2 {
	//	t.Fatal()
	//}
	//if v := metric.Values()[2]; v.Value != 0.1 {
	//	t.Fatal()
	//}
	//time.Sleep(15 * time.Second)
}
