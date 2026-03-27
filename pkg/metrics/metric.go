package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	MetricsNamespace = "easyjet"

	ResultOK    = "ok"
	ResultError = "err"
)

var (
	// TODO: унести в adapter/handler/internal
	DefaultRegistry = prometheus.NewRegistry()

	// TODO: унести в контроллер метрик
	runTime = promauto.With(DefaultRegistry).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: MetricsNamespace,
		Subsystem: "core",
		Name:      "run_duration",
		Help:      "Время выполнения задачи",
		Buckets:   []float64{0.1, 0.5, 2, 5, 10, 30, 60, 120, 300, 600},
	}, []string{"project_id", "result"})
	startTime = promauto.With(DefaultRegistry).NewGauge(prometheus.GaugeOpts{
		Namespace: MetricsNamespace,
		Name:      "start_timestamp",
		Help:      "Время запуска приложения в секундах",
	})
)

func init() {
	DefaultRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(
			// TODO: убрать ненужные метрики
			collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll),
		),
	)
	startTime.Set(float64(time.Now().Unix()))
}

func ConvertOk(ok bool) string {
	if ok {
		return ResultOK
	}

	return ResultError
}

func ObserveRun(projectID uint, ok bool, d time.Duration) {
	runTime.WithLabelValues(strconv.FormatUint(uint64(projectID), 10), ConvertOk(ok)).Observe(d.Seconds())
}
