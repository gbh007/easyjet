package internal

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g *Generator) withVariables(builder *dashboard.DashboardBuilder) *dashboard.DashboardBuilder {
	builder.WithVariable(
		dashboard.NewDatasourceVariableBuilder(metricVariableName).
			Type(prometheusDatasourceType),
	)

	if g.vlExpr != "" {
		builder.WithVariable(
			dashboard.
				NewDatasourceVariableBuilder(logsVariableName).
				Type(logsVariableTypeVictoriaLogs),
		)
	}

	builder.WithVariable(
		dashboard.NewQueryVariableBuilder(projectIDVariableName).
			Query(dashboard.StringOrMap{
				String: new("label_values(easyjet_core_run_duration_count, project_id)"),
			}).
			Datasource(metricDatasource).
			Multi(true).
			IncludeAll(true).
			Refresh(dashboard.VariableRefreshOnTimeRangeChanged),
	)
	builder.WithVariable(
		dashboard.NewQueryVariableBuilder(instanceVariableName).
			Query(dashboard.StringOrMap{
				String: new("label_values(easyjet_start_timestamp, instance)"),
			}).
			Datasource(metricDatasource).
			Multi(true).
			IncludeAll(true).
			Refresh(dashboard.VariableRefreshOnTimeRangeChanged),
	)

	return builder
}
