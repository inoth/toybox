package metric

// https://github.com/zsais/go-gin-prometheus/blob/master/middleware.go

var (
	Counter      string = "counter"
	CounterVec   string = "counter_vec"
	Gauge        string = "gauge"
	GaugeVec     string = "gauge_vec"
	Histogram    string = "histogram"
	HistogramVec string = "histogram_vec"
	Summary      string = "summary"
	SummaryVec   string = "summary_vec"
)

type Prometheus struct {
	metric []*Metric

	Subsystem string `toml:"subsystem"`
}
