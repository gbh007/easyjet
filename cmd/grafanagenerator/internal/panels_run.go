package internal

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g *Generator) runDurationTS() *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title("Run Duration").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`histogram_quantile(0.95, 
						sum(
							rate(easyjet_core_run_duration_bucket{project_id=~"$%s", instance=~"$%s"}[$__rate_interval])
						) by (le, project_id, result)
					)`,
					projectIDVariableName,
					instanceVariableName,
				)).
				LegendFormat("P95 {{project_id}} - {{result}}").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`histogram_quantile(0.80, 
						sum(
							rate(easyjet_core_run_duration_bucket{project_id=~"$%s", instance=~"$%s"}[$__rate_interval])
						) by (le, project_id, result)
					)`,
					projectIDVariableName,
					instanceVariableName,
				)).
				LegendFormat("P80 {{project_id}} - {{result}}").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(rate(easyjet_core_run_duration_sum{project_id=~"$%s", instance=~"$%s"}[$__rate_interval])) by (project_id, result)
					/ sum(rate(easyjet_core_run_duration_count{project_id=~"$%s", instance=~"$%s"}[$__rate_interval])) by (project_id, result)`,
					projectIDVariableName,
					instanceVariableName,
					projectIDVariableName,
					instanceVariableName,
				)).
				LegendFormat("Avg {{project_id}} - {{result}}").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`easyjet_core_last_run_seconds{project_id=~"$%s", instance=~"$%s"}`,
					projectIDVariableName,
					instanceVariableName,
				)).
				LegendFormat("Last {{project_id}} - {{result}}").
				Datasource(metricDatasource),
		}).
		Unit(units.Seconds).
		Legend(simpleLegend()).
		Datasource(metricDatasource)
}

func (g *Generator) runRPSTS() *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title("Run RPS").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(rate(easyjet_core_run_duration_count{project_id=~"$%s", instance=~"$%s"}[$__rate_interval])) by (project_id, result)`,
					projectIDVariableName,
					instanceVariableName,
				)).
				LegendFormat("RPS {{project_id}} - {{result}}").
				Datasource(metricDatasource),
		}).
		Unit(units.RequestsPerSecond).
		Legend(simpleLegend()).
		Datasource(metricDatasource)
}
