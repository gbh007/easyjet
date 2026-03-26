package internal

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

const (
	metricVariableName    = "metrics"
	projectIDVariableName = "project_id"
	instanceVariableName  = "instance"

	prometheusDatasourceType = "prometheus"
)

var metricDatasource = dashboard.DataSourceRef{
	Type: new(prometheusDatasourceType),
	Uid:  new("${" + metricVariableName + "}"),
}
