package internal

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g *Generator) withPanels(builder *dashboard.DashboardBuilder) *dashboard.DashboardBuilder {
	builder.WithRow(dashboard.NewRowBuilder("Resource usage"))
	builder.WithPanel(g.cpuUsageTS())
	builder.WithPanel(g.memUsageTS())
	builder.WithPanel(g.threadsTS())
	builder.WithRow(dashboard.NewRowBuilder("Task Execution"))
	builder.WithPanel(g.runDurationTS())
	builder.WithPanel(g.runRPSTS())
	return builder
}

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
		Datasource(metricDatasource)
}

func (g *Generator) cpuUsageTS() *timeseries.PanelBuilder {
	query := fmt.Sprintf(
		`sum(rate(process_cpu_seconds_total{instance=~"$%s"}[$__rate_interval])) by (instance)`,
		instanceVariableName,
	)
	query3 := fmt.Sprintf(
		`sum(rate(go_cpu_classes_user_cpu_seconds_total{instance=~"$%s"}[$__rate_interval])) by (instance)`,
		instanceVariableName,
	)

	return timeseries.NewPanelBuilder().
		Title("CPU usage").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}} process").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}} process %").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(query3).
				LegendFormat("{{instance}} go code %").
				Datasource(metricDatasource),
		}).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.PercentUnit,
			},
		}).
		OverrideByQuery("C", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.PercentUnit,
			},
		}).
		Unit(units.Seconds).
		Datasource(metricDatasource)
}

func (g *Generator) memUsageTS() *timeseries.PanelBuilder {
	query := fmt.Sprintf(
		`sum(process_resident_memory_bytes{instance=~"$%s"}) by (instance)`,
		instanceVariableName,
	)
	query2 := fmt.Sprintf(
		`sum(go_memstats_sys_bytes{instance=~"$%s"}) by (instance)`,
		instanceVariableName,
	)

	return timeseries.NewPanelBuilder().
		Title("Memory usage").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}} process resident").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(query2).
				LegendFormat("{{instance}} go sys").
				Datasource(metricDatasource),
		}).
		Unit(units.BytesIEC).
		Datasource(metricDatasource)
}

func (g *Generator) threadsTS() *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title("Gorutines & threads").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(go_goroutines{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} gorutines").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(process_num_threads{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} process threads").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(go_threads{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} go threads").
				Datasource(metricDatasource),
		}).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.Number,
			},
		}).
		OverrideByQuery("C", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.Number,
			},
		}).
		Unit(units.Short).
		Datasource(metricDatasource)
}
