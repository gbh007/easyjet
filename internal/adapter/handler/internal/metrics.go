package internal

import "github.com/prometheus/client_golang/prometheus"

var DefaultRegistry = prometheus.NewRegistry()
