package metrics

import (
	"strconv"
	"time"

	"github.com/gbh007/easyjet/internal/adapter/handler/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricsNamespace = "easyjet"

	resultOK    = "ok"
	resultError = "err"
)

var (
	runTime = promauto.With(internal.DefaultRegistry).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Subsystem: "core",
		Name:      "run_duration",
		Help:      "Время выполнения задачи",
		Buckets:   []float64{0.1, 0.5, 2, 5, 10, 30, 60, 120, 300, 600},
	}, []string{"project_id", "result"})
	lastRunTime = promauto.With(internal.DefaultRegistry).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: "core",
		Name:      "last_run_seconds",
		Help:      "Время выполнения последней задачи",
	}, []string{"project_id", "result"})
	startTime = promauto.With(internal.DefaultRegistry).NewGauge(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "start_timestamp",
		Help:      "Время запуска приложения в секундах",
	})
)

func init() {
	internal.DefaultRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(
			// TODO: убрать ненужные метрики
			collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll),
		),
	)
	startTime.Set(float64(time.Now().Unix()))
}

func convertOk(ok bool) string {
	if ok {
		return resultOK
	}

	return resultError
}

func observeRun(projectID uint, ok bool, d time.Duration) {
	runTime.WithLabelValues(strconv.FormatUint(uint64(projectID), 10), convertOk(ok)).Observe(d.Seconds())
	lastRunTime.WithLabelValues(strconv.FormatUint(uint64(projectID), 10), convertOk(ok)).Set(d.Seconds())
}
