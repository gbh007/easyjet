package internal

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

const (
	metricVariableName    = "metrics"
	projectIDVariableName = "project_id"
	instanceVariableName  = "instance"

	prometheusDatasourceType = "prometheus"

	logsVariableTypeVictoriaLogs = "victoriametrics-logs-datasource"
	logsVariableName             = "logs"
)

var (
	metricDatasource = dashboard.DataSourceRef{
		Type: new(prometheusDatasourceType),
		Uid:  new("${" + metricVariableName + "}"),
	}

	logsVictoriaLogsDatasource = dashboard.DataSourceRef{
		Type: new(logsVariableTypeVictoriaLogs),
		Uid:  new("${" + logsVariableName + "}"),
	}
)
