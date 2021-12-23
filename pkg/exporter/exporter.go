package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	calendarEventInfo = prometheus.NewDesc(prometheus.BuildFQName("calendar", "event", "info"), "Info on a calendar event", []string{"id", "location", "description", "summary", "start_time", "end_time"}, nil)
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- calendarEventInfo
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(calendarEventInfo, prometheus.GaugeValue, 1, "uid", "location", "description", "summary", "start_time", "end_time")
}
